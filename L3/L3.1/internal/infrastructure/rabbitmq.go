package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"L3.1/internal/config"
	"L3.1/internal/models"
)

// RabbitMQClient параметры подключения к rabbitmq
type RabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewRabbitMQClient конструктор для структуры клиента
func NewRabbitMQClient(config config.RabbitMQConfig) (*RabbitMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@localhost:5672/", config.User, config.Password)
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish отпрявляет сообщение в rabbitmq и создает очередь если ее нет
func (rc *RabbitMQClient) Publish(ctx context.Context, queueName string, message *models.RabbitMQMessage) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// durable=true → очередь сохранится после рестарта RabbitMQ
	_, err = rc.channel.QueueDeclare(
		queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = rc.channel.PublishWithContext(ctx,
		"", queueName, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Consume — читает сообщения из очереди и вызывает обработчик для каждого.
func (rc *RabbitMQClient) Consume(ctx context.Context, queueName string, handler func(msg *models.RabbitMQMessage) error) error {
	msgs, err := rc.channel.Consume(
		queueName,
		"", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for d := range msgs {
			var m models.RabbitMQMessage
			if err := json.Unmarshal(d.Body, &m); err != nil {
				d.Nack(false, false) // отклоняем без повторной доставки
				continue
			}

			if err := handler(&m); err != nil {
				// повторная доставка при ошибке
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	// блокируем до завершения контекста
	<-ctx.Done()
	return ctx.Err()
}

// RetryMessage — повторная отправка сообщения через delay.
func (rc *RabbitMQClient) RetryMessage(ctx context.Context, queueName string, message *models.RabbitMQMessage, delaySeconds int) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	message.RetryCount++
	return rc.Publish(ctx, queueName, message)
}

// QueueLength — возвращает количество сообщений в очереди.
func (rc *RabbitMQClient) QueueLength(ctx context.Context, queueName string) (int, error) {
	q, err := rc.channel.QueueInspect(queueName)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect queue: %w", err)
	}
	return q.Messages, nil
}
