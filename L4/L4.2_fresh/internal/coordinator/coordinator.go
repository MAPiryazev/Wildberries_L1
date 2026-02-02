package coordinator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/broker"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
	"github.com/rs/zerolog"
)

type ChunkTask struct {
	ID        string
	Content   string
	Checksum  string
	Delimiter string
	Fields    []int
	Suppress  bool
}

type ChunkResult struct {
	TaskID   string
	Output   string
	WorkerID string
	Error    string
}

type Coordinator struct {
	processor      cut.Processor
	quorumSize     int
	logger         zerolog.Logger
	broker         *broker.Broker
	mu             sync.RWMutex
	totalTasks     atomic.Uint64
	completedTasks atomic.Uint64
}

func NewCoordinator(
	processor cut.Processor,
	quorumSize int,
	b *broker.Broker,
	logger zerolog.Logger,
) *Coordinator {
	return &Coordinator{
		processor:  processor,
		quorumSize: quorumSize,
		logger:     logger,
		broker:     b,
	}
}

func (c *Coordinator) SplitIntoChunks(r io.Reader, chunkSize int) ([]ChunkTask, error) {
	var chunks []ChunkTask
	scanner := bufio.NewScanner(r)
	buffer := ""
	bufferSize := 0
	chunkID := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineSize := len(line) + 1

		if bufferSize+lineSize > chunkSize && bufferSize > 0 {
			chunks = append(chunks, ChunkTask{
				ID:      fmt.Sprintf("chunk-%d", chunkID),
				Content: buffer,
			})
			buffer = line + "\n"
			bufferSize = lineSize
			chunkID++
		} else {
			buffer += line + "\n"
			bufferSize += lineSize
		}
	}

	if bufferSize > 0 {
		chunks = append(chunks, ChunkTask{
			ID:      fmt.Sprintf("chunk-%d", chunkID),
			Content: buffer,
		})
	}

	return chunks, scanner.Err()
}

func (c *Coordinator) PublishTasks(
	ctx context.Context,
	chunks []ChunkTask,
	delimiter string,
	fields []int,
	suppress bool,
) error {
	for i := range chunks {
		chunks[i].Delimiter = delimiter
		chunks[i].Fields = fields
		chunks[i].Suppress = suppress

		task := &broker.TaskMessage{
			ID:        chunks[i].ID,
			Chunk:     chunks[i].Content,
			Delimiter: delimiter,
			Fields:    fields,
			Suppress:  suppress,
		}

		if err := c.broker.PublishTask(ctx, task); err != nil {
			return fmt.Errorf("failed to publish task %s: %w", chunks[i].ID, err)
		}
	}

	c.totalTasks.Store(uint64(len(chunks)))
	c.logger.Info().
		Int("count", len(chunks)).
		Msg("published all tasks")

	return nil
}

func (c *Coordinator) CollectResults(
	ctx context.Context,
	expectedCount int,
	timeout time.Duration,
) (map[string]*ChunkResult, error) {
	resultsCh, err := c.broker.ConsumeResult(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to consume results: %w", err)
	}

	results := make(map[string]*ChunkResult)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	collected := 0

	for {
		select {
		case <-timeoutCtx.Done():
			return results, fmt.Errorf(
				"timeout collecting results: got %d/%d, error: %w",
				collected,
				expectedCount,
				timeoutCtx.Err(),
			)
		case result, ok := <-resultsCh:
			if !ok {
				return results, fmt.Errorf("results channel closed unexpectedly")
			}

			results[result.TaskID] = &ChunkResult{
				TaskID:   result.TaskID,
				Output:   result.Output,
				WorkerID: result.WorkerID,
				Error:    result.Error,
			}

			collected++
			c.completedTasks.Add(1)

			c.logger.Debug().
				Str("task_id", result.TaskID).
				Str("worker_id", result.WorkerID).
				Int("completed", collected).
				Int("expected", expectedCount).
				Msg("result received")

			if collected >= expectedCount {
				return results, nil
			}
		}
	}
}

func (c *Coordinator) CheckQuorum(
	results map[string]*ChunkResult,
	expectedCount int,
) (bool, error) {
	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.Error != "" {
			errorCount++
		} else {
			successCount++
		}
	}

	threshold := (expectedCount / 2) + 1

	c.logger.Info().
		Int("success", successCount).
		Int("error", errorCount).
		Int("threshold", threshold).
		Msg("quorum check")

	if successCount >= threshold {
		return true, nil
	}

	return false, fmt.Errorf(
		"quorum not reached: %d/%d workers succeeded, need %d",
		successCount,
		expectedCount,
		threshold,
	)
}

func (c *Coordinator) ProcessWithQuorum(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	delimiter string,
	fields []int,
	suppress bool,
	chunkSize int,
	timeout time.Duration,
) error {
	chunks, err := c.SplitIntoChunks(r, chunkSize)
	if err != nil {
		return fmt.Errorf("failed to split into chunks: %w", err)
	}

	if len(chunks) == 0 {
		return fmt.Errorf("no chunks to process")
	}

	c.logger.Info().
		Int("chunks", len(chunks)).
		Int("chunk_size", chunkSize).
		Msg("starting distributed processing")

	if err := c.PublishTasks(ctx, chunks, delimiter, fields, suppress); err != nil {
		return fmt.Errorf("failed to publish tasks: %w", err)
	}

	results, err := c.CollectResults(ctx, len(chunks), timeout)
	if err != nil {
		c.logger.Warn().Err(err).Msg("failed to collect all results")
	}

	if len(results) == 0 {
		return fmt.Errorf("no results received from workers")
	}

	ok, err := c.CheckQuorum(results, len(chunks))
	if !ok {
		return fmt.Errorf("quorum check failed: %w", err)
	}

	for _, chunk := range chunks {
		result, ok := results[chunk.ID]
		if !ok || result.Error != "" {
			c.logger.Warn().
				Str("chunk_id", chunk.ID).
				Msg("chunk failed, skipping")
			continue
		}

		if _, err := fmt.Fprint(w, result.Output); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	c.logger.Info().
		Uint64("total_tasks", c.totalTasks.Load()).
		Uint64("completed_tasks", c.completedTasks.Load()).
		Msg("distributed processing completed")

	return nil
}

func (c *Coordinator) Stats() (uint64, uint64) {
	return c.totalTasks.Load(), c.completedTasks.Load()
}
