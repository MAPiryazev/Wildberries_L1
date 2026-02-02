package config

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	ErrKafkaParamNotFound = errors.New(`один или несколько критических параметров Kafka не были найдены в env, проверьте:
	KAFKA_BROKERS
	KAFKA_TOPIC_IMAGE_TASKS
	KAFKA_GROUP_ID`)
)

type KafkaConfig struct {
	Brokers             string
	TopicImageTasks     string
	RetryTopic          string
	DeadLetterTopic     string
	GroupID             string
	AutoOffsetReset     string
	EnableAutoCommit    bool
	EnableIdempotence   bool
	IsolationLevel      string
	SessionTimeoutMS    int
	HeartbeatIntervalMS int
	MaxPollIntervalMS   int
}

func LoadKafkaConfig(envPath string) (*KafkaConfig, error) {
	if err := loadEnvFiles(envPath); err != nil {
		return nil, err
	}

	enableAutoCommit, err := parseBoolEnv("KAFKA_ENABLE_AUTO_COMMIT")
	if err != nil {
		log.Println("ошибка при считывании KAFKA_ENABLE_AUTO_COMMIT из env:", err)
	}
	enableIdempotence, err := parseBoolEnv("KAFKA_ENABLE_IDEMPOTENCE")
	if err != nil {
		log.Println("ошибка при считывании KAFKA_ENABLE_IDEMPOTENCE из env:", err)
	}

	sessionTimeout, err := parseIntEnv("KAFKA_SESSION_TIMEOUT_MS")
	if err != nil {
		log.Println("ошибка при считывании KAFKA_SESSION_TIMEOUT_MS из env:", err)
	}
	heartbeatInterval, err := parseIntEnv("KAFKA_HEARTBEAT_INTERVAL_MS")
	if err != nil {
		log.Println("ошибка при считывании KAFKA_HEARTBEAT_INTERVAL_MS из env:", err)
	}
	maxPollInterval, err := parseIntEnv("KAFKA_MAX_POLL_INTERVAL_MS")
	if err != nil {
		log.Println("ошибка при считывании KAFKA_MAX_POLL_INTERVAL_MS из env:", err)
	}

	cfg := &KafkaConfig{
		Brokers:             normalizeKafkaBrokers(os.Getenv("KAFKA_BROKERS")),
		TopicImageTasks:     os.Getenv("KAFKA_TOPIC_IMAGE_TASKS"),
		RetryTopic:          os.Getenv("KAFKA_RETRY_TOPIC"),
		DeadLetterTopic:     os.Getenv("KAFKA_DEAD_LETTER_TOPIC"),
		GroupID:             os.Getenv("KAFKA_GROUP_ID"),
		AutoOffsetReset:     os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
		EnableAutoCommit:    enableAutoCommit,
		EnableIdempotence:   enableIdempotence,
		IsolationLevel:      os.Getenv("KAFKA_ISOLATION_LEVEL"),
		SessionTimeoutMS:    sessionTimeout,
		HeartbeatIntervalMS: heartbeatInterval,
		MaxPollIntervalMS:   maxPollInterval,
	}

	if err := validateKafkaConfig(cfg); err != nil {
		return nil, err
	}

	if strings.TrimSpace(cfg.RetryTopic) == "" {
		log.Println("предупреждение: KAFKA_RETRY_TOPIC не указан")
	}
	if strings.TrimSpace(cfg.DeadLetterTopic) == "" {
		log.Println("предупреждение: KAFKA_DEAD_LETTER_TOPIC не указан")
	}

	return cfg, nil
}

func parseIntEnv(key string) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("некорректное числовое значение %q", value)
	}

	return parsed, nil
}

func validateKafkaConfig(cfg *KafkaConfig) error {
	if cfg.Brokers == "" ||
		cfg.TopicImageTasks == "" ||
		cfg.GroupID == "" {
		return ErrKafkaParamNotFound
	}

	return nil
}

func normalizeKafkaBrokers(brokers string) string {
	if strings.TrimSpace(brokers) == "" || runningInsideDocker() {
		return brokers
	}

	parts := strings.Split(brokers, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		host, port, err := net.SplitHostPort(part)
		if err != nil || host == "" {
			parts[i] = part
			continue
		}
		if host == "kafka" {
			parts[i] = net.JoinHostPort("localhost", port)
		} else {
			parts[i] = net.JoinHostPort(host, port)
		}
	}

	return strings.Join(parts, ",")
}
