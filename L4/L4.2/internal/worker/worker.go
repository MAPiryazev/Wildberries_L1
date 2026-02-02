package worker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/broker"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
	"github.com/rs/zerolog"
)

type Worker struct {
	id        string
	broker    *broker.Broker
	processor cut.Processor
	logger    zerolog.Logger

	threads      int
	processed    atomic.Uint64
	errors       atomic.Uint64
	mu           sync.Mutex
	taskHandlers map[string]context.CancelFunc
}

func NewWorker(
	id string,
	b *broker.Broker,
	processor cut.Processor,
	threads int,
	logger zerolog.Logger,
) *Worker {
	return &Worker{
		id:           id,
		broker:       b,
		processor:    processor,
		logger:       logger,
		threads:      threads,
		taskHandlers: make(map[string]context.CancelFunc),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info().
		Str("worker_id", w.id).
		Int("threads", w.threads).
		Msg("worker started")

	tasksCh, err := w.broker.ConsumeTask(ctx, w.threads)
	if err != nil {
		return fmt.Errorf("failed to consume tasks: %w", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < w.threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			w.processLoop(ctx, threadID, tasksCh)
		}(i)
	}

	wg.Wait()
	w.logger.Info().
		Uint64("processed", w.processed.Load()).
		Uint64("errors", w.errors.Load()).
		Msg("worker stopped")

	return nil
}

func (w *Worker) processLoop(ctx context.Context, threadID int, tasksCh <-chan *broker.TaskMessage) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasksCh:
			if !ok {
				return
			}

			if err := w.handleTask(ctx, task); err != nil {
				w.errors.Add(1)
				w.logger.Error().
					Err(err).
					Str("task_id", task.ID).
					Int("thread_id", threadID).
					Msg("task processing error")
			} else {
				w.processed.Add(1)
			}
		}
	}
}

func (w *Worker) handleTask(ctx context.Context, task *broker.TaskMessage) error {
	logger := w.logger.With().
		Str("task_id", task.ID).
		Str("worker_id", w.id).
		Logger()

	start := time.Now()
	defer func() {
		logger.Info().
			Dur("elapsed", time.Since(start)).
			Msg("task completed")

	}()

	lines := strings.Split(strings.TrimSpace(task.Chunk), "\n")
	var results []string

	for _, line := range lines {
		if line == "" {
			continue
		}

		output, err := w.processor.ProcessLine(line)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("line", line).
				Msg("failed to process line")
			continue
		}

		results = append(results, output)
	}

	output := strings.Join(results, "\n")
	if len(results) > 0 {
		output += "\n"
	}

	result := &broker.ResultMessage{
		TaskID:   task.ID,
		Output:   output,
		WorkerID: w.id,
	}

	if err := w.broker.PublishResult(ctx, result); err != nil {
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for taskID, cancel := range w.taskHandlers {
		cancel()
		w.logger.Debug().Str("task_id", taskID).Msg("cancelled task")
	}
	w.taskHandlers = make(map[string]context.CancelFunc)
}
