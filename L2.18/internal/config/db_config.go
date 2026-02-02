package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

// LoadDBPSQLConfig возвращает параметры подключения к PSQL БД
func LoadDBPSQLConfig(path string) (*DBPSQLConfig, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки env файла с путем %v: %w", path, err)
	}

	tempHost := os.Getenv("POSTGRES_HOST")
	tempPort := os.Getenv("POSTGRES_PORT")
	tempUser := os.Getenv("POSTGRES_USER")
	tempName := os.Getenv("POSTGRES_DB")
	tempPassword := os.Getenv("POSTGRES_PASSWORD")
	tempSSLMode := os.Getenv("POSTGRES_SSLMODE")

	if tempHost == "" || tempPort == "" || tempUser == "" || tempPassword == "" || tempSSLMode == "" || tempName == "" {
		return nil, fmt.Errorf("не найдены критически важные параметры подключения БД")
	}

	tempMaxOpenConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	if err != nil {
		tempMaxOpenConns = 25
	}
	tempMaxIDLEConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	if err != nil {
		tempMaxIDLEConns = 10
	}
	tempMaxConnLifetime, err := strconv.Atoi(os.Getenv("POSTGRES_CONN_MAX_LIFETIME"))
	if err != nil {
		tempMaxConnLifetime = 15
	}

	return &DBPSQLConfig{
		Host:            tempHost,
		Port:            tempPort,
		User:            tempUser,
		DBName:          tempName,
		Password:        tempPassword,
		SSLMode:         tempSSLMode,
		MaxOpenConns:    tempMaxOpenConns,
		MaxIdleConns:    tempMaxIDLEConns,
		MaxConnLifeTime: tempMaxConnLifetime,
	}, nil

}
