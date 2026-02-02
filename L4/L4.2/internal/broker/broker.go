package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type TaskMessage struct {
	ID        string `json:"id"`
	Chunk     string `json:"chunk"`
	Delimiter string `json:"delimiter"`
	Fields    []int  `json:"fields"`
	Suppress  bool   `json:"suppress"`
}

type ResultMessage struct {
	TaskID   string `json:"task_id"`
	Output   string `json:"output"`
	WorkerID string `json:"worker_id"`
	Error    string `json:"error,omitempty"`
}

type Broker struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	taskQueue   amqp.Queue
	resultQueue amqp.Queue
	logger      zerolog.Logger
	mu          sync.RWMutex
	closed      bool
}

func NewBroker(url string, logger zerolog.Logger) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	taskQueue, err := ch.QueueDeclare(
		"mycut.tasks",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare task queue: %w", err)
	}

	resultQueue, err := ch.QueueDeclare(
		"mycut.results",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare result queue: %w", err)
	}

	return &Broker{
		conn:        conn,
		channel:     ch,
		taskQueue:   taskQueue,
		resultQueue: resultQueue,
		logger:      logger,
	}, nil
}

func (b *Broker) PublishTask(ctx context.Context, task *TaskMessage) error {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return fmt.Errorf("broker is closed")
	}
	b.mu.RUnlock()

	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	return b.channel.PublishWithContext(
		ctx,
		"",
		"mycut.tasks",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (b *Broker) ConsumeTask(ctx context.Context, prefetch int) (<-chan *TaskMessage, error) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return nil, fmt.Errorf("broker is closed")
	}
	b.mu.RUnlock()

	if err := b.channel.Qos(prefetch, 0, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := b.channel.ConsumeWithContext(
		ctx,
		"mycut.tasks",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	taskCh := make(chan *TaskMessage, prefetch)
	go func() {
		defer close(taskCh)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var task TaskMessage
				if err := json.Unmarshal(msg.Body, &task); err != nil {
					b.logger.Error().Err(err).Msg("failed to unmarshal task")
					msg.Nack(false, true)
					continue
				}

				select {
				case taskCh <- &task:
					msg.Ack(false)
				case <-ctx.Done():
					msg.Nack(false, true)
					return
				}
			}
		}
	}()

	return taskCh, nil
}

func (b *Broker) PublishResult(ctx context.Context, result *ResultMessage) error {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return fmt.Errorf("broker is closed")
	}
	b.mu.RUnlock()

	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	return b.channel.PublishWithContext(
		ctx,
		"",
		"mycut.results",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (b *Broker) ConsumeResult(ctx context.Context, prefetch int) (<-chan *ResultMessage, error) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return nil, fmt.Errorf("broker is closed")
	}
	b.mu.RUnlock()

	if err := b.channel.Qos(prefetch, 0, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := b.channel.ConsumeWithContext(
		ctx,
		"mycut.results",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	resultCh := make(chan *ResultMessage, prefetch)
	go func() {
		defer close(resultCh)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var result ResultMessage
				if err := json.Unmarshal(msg.Body, &result); err != nil {
					b.logger.Error().Err(err).Msg("failed to unmarshal result")
					msg.Nack(false, true)
					continue
				}

				select {
				case resultCh <- &result:
					msg.Ack(false)
				case <-ctx.Done():
					msg.Nack(false, true)
					return
				}
			}
		}
	}()

	return resultCh, nil
}

func (b *Broker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	b.closed = true

	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
	return nil
}

func (b *Broker) IsClosed() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.closed
}
