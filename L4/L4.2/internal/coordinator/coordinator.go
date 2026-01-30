package coordinator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
	"github.com/rs/zerolog"
)

type ChunkTask struct {
	ID       string
	Content  string
	Checksum string
}

type ChunkResult struct {
	TaskID   string
	Output   string
	WorkerID string
	Error    string
}

type Coordinator struct {
	processor  cut.Processor
	quorumSize int
	logger     zerolog.Logger
	mu         sync.RWMutex
}

func NewCoordinator(processor cut.Processor, quorumSize int, logger zerolog.Logger) *Coordinator {
	return &Coordinator{
		processor:  processor,
		quorumSize: quorumSize,
		logger:     logger,
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

func (c *Coordinator) CheckQuorum(results []ChunkResult) (bool, string, error) {
	if len(results) < c.quorumSize {
		return false, "", fmt.Errorf("insufficient results: got %d, need %d", len(results), c.quorumSize)
	}

	resultMap := make(map[string]int)
	for _, r := range results {
		if r.Error != "" {
			continue
		}
		resultMap[r.Output]++
	}

	if len(resultMap) == 0 {
		return false, "", fmt.Errorf("all workers returned errors")
	}

	var maxOutput string
	maxCount := 0
	for output, count := range resultMap {
		if count > maxCount {
			maxCount = count
			maxOutput = output
		}
	}

	if maxCount >= (c.quorumSize/2 + 1) {
		return true, maxOutput, nil
	}

	return false, "", fmt.Errorf("quorum not reached: max agreement %d/%d", maxCount, c.quorumSize)
}

func (c *Coordinator) ProcessWithQuorum(ctx context.Context, r io.Reader, w io.Writer) error {
	chunks, err := c.SplitIntoChunks(r, 1024*1024)
	if err != nil {
		return err
	}

	c.logger.Info().Int("chunks", len(chunks)).Msg("split input into chunks")

	for _, chunk := range chunks {
		output, err := c.processor.ProcessLine(chunk.Content)
		if err != nil {
			c.logger.Error().Err(err).Str("chunk_id", chunk.ID).Msg("failed to process chunk")
			continue
		}

		if _, err := fmt.Fprintln(w, output); err != nil {
			return err
		}
	}

	return nil
}
