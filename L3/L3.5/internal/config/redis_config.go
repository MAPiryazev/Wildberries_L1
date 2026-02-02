package config

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

// RedisConfig параметры подключения к Redis
type RedisConfig struct {
	Port string
}

// LoadRedisConfig возвращает параметры подключения к Redis
func LoadRedisConfig(path string) (*RedisConfig, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("путь до env файла не был передан ")
	}

	cfg := config.New()
	err := cfg.LoadEnvFiles(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке .env файла %w", err)
	}
	cfg.EnableEnv("")

	tempRedisConfig := &RedisConfig{}
	tempRedisConfig.Port = cfg.GetString("redis.port")

	if len(tempRedisConfig.Port) == 0 {
		return nil, fmt.Errorf("не был найден порт для подключения к redis")
	}
	return tempRedisConfig, nil
}
