package config

import (
	"fmt"
	"time"

	wbfconfig "github.com/wb-go/wbf/config"
)

type Config struct {
	Database Database
	Server   Server
	App      App
	JWT      JWT
}

type Database struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
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

type JWT struct {
	Secret     string
	Expiration time.Duration
}

func Load(envPath ...string) (*Config, error) {
	cfg := wbfconfig.New()

	path := "environment"
	if len(envPath) > 0 && envPath[0] != "" {
		path = envPath[0]
	}

	if err := cfg.LoadEnvFiles(path); err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	cfg.EnableEnv("")

	cfg.SetDefault("postgres.host", "localhost")
	cfg.SetDefault("postgres.port", 5433)
	cfg.SetDefault("postgres.user", "warehouse_user")
	cfg.SetDefault("postgres.password", "warehouse_secret")
	cfg.SetDefault("postgres.name", "warehouse_control")
	cfg.SetDefault("postgres.max_open_conns", 25)
	cfg.SetDefault("postgres.max_idle_conns", 10)
	cfg.SetDefault("postgres.conn_max_lifetime", "15m")

	cfg.SetDefault("server.port", 8080)
	cfg.SetDefault("server.read_timeout", 10)
	cfg.SetDefault("server.write_timeout", 10)

	cfg.SetDefault("app.env", "development")
	cfg.SetDefault("app.log_level", "debug")

	cfg.SetDefault("jwt.secret", "jwt_secret")
	cfg.SetDefault("jwt.expiration", "24h")

	lifetime, err := time.ParseDuration(cfg.GetString("postgres.conn_max_lifetime"))
	if err != nil {
		lifetime = 15 * time.Minute
	}

	expiration, err := time.ParseDuration(cfg.GetString("jwt.expiration"))
	if err != nil {
		expiration = 24 * time.Hour
	}

	appCfg := &Config{
		Database: Database{
			Host:            cfg.GetString("postgres.host"),
			Port:            cfg.GetInt("postgres.port"),
			User:            cfg.GetString("postgres.user"),
			Password:        cfg.GetString("postgres.password"),
			Name:            cfg.GetString("postgres.name"),
			MaxOpenConns:    cfg.GetInt("postgres.max_open_conns"),
			MaxIdleConns:    cfg.GetInt("postgres.max_idle_conns"),
			ConnMaxLifetime: lifetime,
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
		JWT: JWT{
			Secret:     cfg.GetString("jwt.secret"),
			Expiration: expiration,
		},
	}

	return appCfg, nil
}
