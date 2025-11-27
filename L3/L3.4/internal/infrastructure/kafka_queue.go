package infrastructure

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/config"
)

// KafkaQueue реализация MessageQueue интерфейса для Kafka
type KafkaQueue struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	config   *config.KafkaConfig
	wg       sync.WaitGroup
	cancel   context.CancelFunc
}

// NewKafkaQueue создает новый клиент Kafka Queue
func NewKafkaQueue(cfg *config.KafkaConfig) (*KafkaQueue, error) {
	log.Printf("[kafka] init queue, brokers=%s topic=%s group=%s idempotent=%v",
		cfg.Brokers, cfg.TopicImageTasks, cfg.GroupID, cfg.EnableIdempotence)
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Idempotent = cfg.EnableIdempotence
	if cfg.EnableIdempotence {
		producerConfig.Net.MaxOpenRequests = 1
	}

	brokers := []string{cfg.Brokers}
	if len(cfg.Brokers) > 0 {
		brokerList := []string{}
		for _, broker := range splitByComma(cfg.Brokers) {
			if len(broker) > 0 {
				brokerList = append(brokerList, broker)
			}
		}
		if len(brokerList) > 0 {
			brokers = brokerList
		}
	}

	producer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka producer: %w", err)
	}
	log.Printf("[kafka] producer connected (%d brokers)", len(brokers))

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	if cfg.AutoOffsetReset == "latest" {
		consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	consumerConfig.Consumer.Group.Session.Timeout = time.Duration(cfg.SessionTimeoutMS) * time.Millisecond
	consumerConfig.Consumer.Group.Heartbeat.Interval = time.Duration(cfg.HeartbeatIntervalMS) * time.Millisecond
	consumerConfig.Consumer.MaxProcessingTime = time.Duration(cfg.MaxPollIntervalMS) * time.Millisecond

	consumer, err := sarama.NewConsumerGroup(brokers, cfg.GroupID, consumerConfig)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("ошибка создания Kafka consumer: %w", err)
	}
	log.Printf("[kafka] consumer group connected")

	return &KafkaQueue{
		producer: producer,
		consumer: consumer,
		config:   cfg,
	}, nil
}

// Publish публикует сообщение в Kafka topic
func (k *KafkaQueue) Publish(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("ошибка публикации сообщения в topic %s: %w", topic, err)
	}

	log.Printf("Сообщение опубликовано в topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

// Subscribe подписывается на Kafka topic и обрабатывает сообщения
func (k *KafkaQueue) Subscribe(topic, groupID string, handler func(msg []byte) error) error {
	ctx, cancel := context.WithCancel(context.Background())
	k.cancel = cancel

	handlerWrapper := &kafkaConsumerGroupHandler{
		handler: handler,
		topic:   topic,
	}

	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := k.consumer.Consume(ctx, []string{topic}, handlerWrapper); err != nil {
					log.Printf("Ошибка при получении сообщений из topic %s: %v", topic, err)
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	return nil
}

// Close закрывает соединения с Kafka
func (k *KafkaQueue) Close() error {
	if k.cancel != nil {
		k.cancel()
	}

	k.wg.Wait()

	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия producer: %w", err)
	}

	if err := k.consumer.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия consumer: %w", err)
	}

	return nil
}

// kafkaConsumerGroupHandler реализует sarama.ConsumerGroupHandler
type kafkaConsumerGroupHandler struct {
	handler func(msg []byte) error
	topic   string
}

func (h *kafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := h.handler(message.Value); err != nil {
				log.Printf("Ошибка обработки сообщения из topic %s: %v", h.topic, err)
				// Можно добавить логику retry или отправки в dead letter queue
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

// splitByComma разделяет строку по запятым
func splitByComma(s string) []string {
	result := strings.Split(s, ",")
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
	}
	return result
}
