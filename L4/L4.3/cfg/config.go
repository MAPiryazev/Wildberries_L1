package cfg

import (
	"fmt"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App  `yaml:"app"`
	Log  string `env-required:"true" yaml:"log" env:"LOG_LEVEL"`
	Host string `env-required:"true" yaml:"host" env:"HOST"`
	Port string `env-required:"true" yaml:"port" env:"PORT"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

// NewConfig Configuration init
func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	_ = godotenv.Load()
	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
