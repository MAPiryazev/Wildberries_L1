package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel          string
	Delimiter         string
	Fields            string
	SuppressNoDelim   bool
	RabbitMQURL       string
	WorkerID          string
	CoordinatorListen string
	WorkerThreads     int
	ChunkSize         int
	QuorumSize        int
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
		LogLevel:          viper.GetString("LOG_LEVEL"),
		Delimiter:         viper.GetString("DELIMITER"),
		Fields:            viper.GetString("FIELDS"),
		SuppressNoDelim:   viper.GetBool("SUPPRESS_NO_DELIM"),
		RabbitMQURL:       viper.GetString("RABBITMQ_URL"),
		WorkerID:          viper.GetString("WORKER_ID"),
		CoordinatorListen: viper.GetString("COORDINATOR_LISTEN"),
		WorkerThreads:     viper.GetInt("WORKER_THREADS"),
		ChunkSize:         viper.GetInt("CHUNK_SIZE"),
		QuorumSize:        viper.GetInt("QUORUM_SIZE"),
	}

	return nil
}

func BindFlags(flagSet *pflag.FlagSet) {
	viper.BindPFlag("DELIMITER", flagSet.Lookup("d"))
	viper.BindPFlag("FIELDS", flagSet.Lookup("f"))
	viper.BindPFlag("SUPPRESS_NO_DELIM", flagSet.Lookup("s"))
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
