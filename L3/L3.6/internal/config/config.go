package config

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Database Database
	Server   Server
	App      App
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Server struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type App struct {
	Env      string
	LogLevel string
}

func Load() (*Config, error) {
	cfg := config.New()

	if err := cfg.LoadEnvFiles("environment"); err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	cfg.EnableEnv("")

	cfg.SetDefault("db.host", "localhost")
	cfg.SetDefault("db.port", 5432)
	cfg.SetDefault("db.user", "postgres")
	cfg.SetDefault("db.password", "password")
	cfg.SetDefault("db.name", "salestracker")
	cfg.SetDefault("server.port", 8080)
	cfg.SetDefault("server.read_timeout", 10)
	cfg.SetDefault("server.write_timeout", 10)
	cfg.SetDefault("app.env", "development")
	cfg.SetDefault("app.log_level", "debug")

	appCfg := &Config{
		Database: Database{
			Host:     cfg.GetString("db.host"),
			Port:     cfg.GetInt("db.port"),
			User:     cfg.GetString("db.user"),
			Password: cfg.GetString("db.password"),
			Name:     cfg.GetString("db.name"),
		},
		Server: Server{
			Port:         cfg.GetInt("server.port"),
			ReadTimeout:  cfg.GetInt("server.read_timeout"),
			WriteTimeout: cfg.GetInt("server.write_timeout"),
		},
		App: App{
			Env:      cfg.GetString("app.env"),
			LogLevel: cfg.GetString("app.log_level"),
		},
	}

	return appCfg, nil
}
