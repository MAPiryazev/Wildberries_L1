package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel             string
	Delimiter            string
	Fields               string
	SuppressNoDelim      bool
	RabbitMQURL          string
	RabbitMQTasksQueue   string
	RabbitMQResultsQueue string
	RabbitMQPrefetch     int
	WorkerID             string
	CoordinatorListen    string
	WorkerThreads        int
	ChunkSize            int
	QuorumSize           int
	TimeoutWorker        time.Duration
	TimeoutQuorum        time.Duration
	Mode                 string
}

var cfg *Config

func Init(envFile string) error {
	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	cfg = &Config{
		LogLevel:             viper.GetString("LOG_LEVEL"),
		Delimiter:            viper.GetString("DELIMITER"),
		Fields:               viper.GetString("FIELDS"),
		SuppressNoDelim:      viper.GetBool("SUPPRESS_NO_DELIM"),
		RabbitMQURL:          viper.GetString("RABBITMQ_URL"),
		RabbitMQTasksQueue:   viper.GetString("RABBITMQ_QUEUE_TASKS"),
		RabbitMQResultsQueue: viper.GetString("RABBITMQ_QUEUE_RESULTS"),
		RabbitMQPrefetch:     viper.GetInt("RABBITMQ_PREFETCH"),
		WorkerID:             viper.GetString("WORKER_ID"),
		CoordinatorListen:    viper.GetString("COORDINATOR_LISTEN"),
		WorkerThreads:        viper.GetInt("WORKER_THREADS"),
		ChunkSize:            viper.GetInt("CHUNK_SIZE"),
		QuorumSize:           viper.GetInt("QUORUM_SIZE"),
		TimeoutWorker:        viper.GetDuration("TIMEOUT_WORKER"),
		TimeoutQuorum:        viper.GetDuration("TIMEOUT_QUORUM"),
	}

	if cfg.WorkerThreads <= 0 {
		cfg.WorkerThreads = 4
	}
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = 1024 * 1024
	}
	if cfg.QuorumSize <= 0 {
		cfg.QuorumSize = 3
	}
	if cfg.RabbitMQPrefetch <= 0 {
		cfg.RabbitMQPrefetch = 10
	}
	if cfg.TimeoutWorker <= 0 {
		cfg.TimeoutWorker = 30 * time.Second
	}
	if cfg.TimeoutQuorum <= 0 {
		cfg.TimeoutQuorum = 60 * time.Second
	}

	return nil
}

func BindFlags(flagSet *pflag.FlagSet) {
	viper.BindPFlag("DELIMITER", flagSet.Lookup("d"))
	viper.BindPFlag("FIELDS", flagSet.Lookup("f"))
	viper.BindPFlag("SUPPRESS_NO_DELIM", flagSet.Lookup("s"))
	viper.BindPFlag("MODE", flagSet.Lookup("mode"))
}

func Get() *Config {
	return cfg
}

func GetLogger() zerolog.Logger {
	if cfg == nil {
		return log.Logger
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	return log.With().Logger().Level(level)
}
