package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// APIConfig структура для настроек API
type APIConfig struct {
	Port string
}

// LoadAPIConfig принимает путь до env файла и возвращает порт для api
func LoadAPIConfig(path string) *APIConfig {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("предупреждение: .env файл не найден по пути %s, используем настройки по умолчанию", path)
	}

	return &APIConfig{
		Port: getEnvHelper("API_PORT", "8080"),
	}
}

func getEnvHelper(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
