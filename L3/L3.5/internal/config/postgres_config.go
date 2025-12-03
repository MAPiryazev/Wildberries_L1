package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/wb-go/wbf/dbpg"
)

// DBPSQLConfig параметры подключения к psql БД
type DBPSQLConfig struct {
	Host            string
	Port            string
	User            string
	DBName          string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifeTime int
}

// LoadDBPSQL возвращает DSN строку и *dbpg.Options для подключения через фреймворк WBF
func LoadDBPSQL(path string) (string, *dbpg.Options, error) {
	err := godotenv.Load(path)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка загрузки env файла: %w", err)
	}

	cfg := &DBPSQLConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		return "", nil, fmt.Errorf("не найдены обязательные параметры подключения к БД")
	}

	cfg.MaxOpenConns, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}

	cfg.MaxIdleConns, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10
	}

	cfg.MaxConnLifeTime, _ = strconv.Atoi(os.Getenv("POSTGRES_CONN_MAX_LIFETIME"))
	if cfg.MaxConnLifeTime == 0 {
		cfg.MaxConnLifeTime = 15
	}

	masterDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// Формируем *dbpg.Options
	opts := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.MaxConnLifeTime) * time.Minute,
	}

	return masterDSN, opts, nil
}
