package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ClickHouseConfig параметры подключения к ClickHouse
type ClickHouseConfig struct {
	Host       string
	HTTPPort   string
	NativePort string
	User       string
	Password   string
	Database   string
	Secure     bool
}

// LoadClickHouseConfig возвращает параметры подключения к ClickHouse
func LoadClickHouseConfig(path string) (*ClickHouseConfig, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("ошибка загрузки env файла с путем %v: %w", path, err)
	}

	host := os.Getenv("CLICKHOUSE_HOST")
	httpPort := os.Getenv("CLICKHOUSE_HTTP_PORT")
	nativePort := os.Getenv("CLICKHOUSE_NATIVE_PORT")
	user := os.Getenv("CLICKHOUSE_USER")
	password := os.Getenv("CLICKHOUSE_PASSWORD")
	database := os.Getenv("CLICKHOUSE_DB")

	if httpPort == "" {
		legacyHTTPPort := os.Getenv("CLICKHOUSE_PORT")
		if legacyHTTPPort != "" {
			httpPort = legacyHTTPPort
		}
	}

	if httpPort == "" {
		httpPort = "8123"
	}
	if nativePort == "" {
		nativePort = "9000"
	}

	secure := false
	if v := os.Getenv("CLICKHOUSE_SECURE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			secure = b
		} else if v == "1" {
			secure = true
		}
	}

	if host == "" || user == "" || database == "" {
		return nil, fmt.Errorf("не найдены критически важные параметры подключения ClickHouse (CLICKHOUSE_HOST/CLICKHOUSE_USER/CLICKHOUSE_DB)")
	}

	return &ClickHouseConfig{
		Host:       host,
		HTTPPort:   httpPort,
		NativePort: nativePort,
		User:       user,
		Password:   password,
		Database:   database,
		Secure:     secure,
	}, nil
}
