package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
	"github.com/rs/zerolog"
)

type TaskPayload struct {
	TaskID    string
	Chunk     string
	Delimiter string
	Fields    []int
	Suppress  bool
}

type ResultPayload struct {
	TaskID   string
	Output   string
	Error    string
	WorkerID string
}

type Worker struct {
	id        string
	processor cut.Processor
	logger    zerolog.Logger
}

func NewWorker(id string, processor cut.Processor, logger zerolog.Logger) *Worker {
	return &Worker{
		id:        id,
		processor: processor,
		logger:    logger,
	}
}

func (w *Worker) ProcessTask(ctx context.Context, payload TaskPayload) (*ResultPayload, error) {
	lines := len(payload.Chunk)
	w.logger.Info().
		Str("task_id", payload.TaskID).
		Int("size", lines).
		Msg("processing task")

	output, err := w.processor.ProcessLine(payload.Chunk)
	if err != nil {
		return &ResultPayload{
			TaskID:   payload.TaskID,
			Error:    err.Error(),
			WorkerID: w.id,
		}, nil
	}

	return &ResultPayload{
		TaskID:   payload.TaskID,
		Output:   output,
		WorkerID: w.id,
	}, nil
}

func (w *Worker) ProcessTaskFromJSON(ctx context.Context, data []byte) (*ResultPayload, error) {
	var task TaskPayload
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return w.ProcessTask(ctx, task)
}
