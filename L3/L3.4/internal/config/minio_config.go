package config

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrEnvNotFound        = errors.New("не удалось загрузить env файл")
	ErrMinioParamNotFound = errors.New(`один или несколько критических параметров не были найдены в env, проверьте:
	MINIO_ENDPOINT
	MINIO_ACCESS_KEY
	MINIO_SECRET_KEY
	MINIO_BUCKET_ORIGINAL
	MINIO_BUCKET_PROCESSED`)
)

type MinioConfgig struct {
	MinioEndpoint        string
	MinioAccessKey       string
	MinioSecretKey       string
	MinioUseSSL          bool
	MinioBucketOriginal  string
	MinioBucketProcessed string
	MinioBucketMetadata  string
}

func LoadMinioConfig(envPath string) (*MinioConfgig, error) {
	if err := loadEnvFiles(envPath); err != nil {
		return nil, err
	}

	useSSL, err := parseBoolEnv("MINIO_USE_SSL")
	if err != nil {
		log.Println("ошибка при считывании MINIO_USE_SSL из env:", err)
	}

	cfg := &MinioConfgig{
		MinioEndpoint:        os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey:       os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey:       os.Getenv("MINIO_SECRET_KEY"),
		MinioUseSSL:          useSSL,
		MinioBucketOriginal:  os.Getenv("MINIO_BUCKET_ORIGINAL"),
		MinioBucketProcessed: os.Getenv("MINIO_BUCKET_PROCESSED"),
		MinioBucketMetadata:  os.Getenv("MINIO_BUCKET_METADATA"),
	}

	if strings.TrimSpace(cfg.MinioBucketMetadata) == "" {
		cfg.MinioBucketMetadata = "metadata"
	}

	cfg.MinioEndpoint = normalizeMinioEndpoint(cfg.MinioEndpoint)

	if err := validateMinioConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadEnvFiles(envPath string) error {
	candidates := []string{}
	if strings.TrimSpace(envPath) != "" {
		candidates = append(candidates, envPath)
	}
	candidates = append(candidates, ".env", "../environment/.env")

	var lastErr error
	for _, path := range candidates {
		if err := godotenv.Load(path); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	return fmt.Errorf("%w: %v", ErrEnvNotFound, lastErr)
}

func parseBoolEnv(key string) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return false, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("некорректное булево значение %q", value)
	}

	return parsed, nil
}

func validateMinioConfig(cfg *MinioConfgig) error {
	if cfg.MinioEndpoint == "" ||
		cfg.MinioAccessKey == "" ||
		cfg.MinioSecretKey == "" ||
		cfg.MinioBucketOriginal == "" ||
		cfg.MinioBucketProcessed == "" {
		return ErrMinioParamNotFound
	}

	return nil
}

func normalizeMinioEndpoint(endpoint string) string {
	if endpoint == "" {
		return endpoint
	}

	host, port, err := net.SplitHostPort(endpoint)
	if err != nil || host == "" || runningInsideDocker() {
		return endpoint
	}

	if host == "minio" {
		return net.JoinHostPort("localhost", port)
	}

	return endpoint
}

func runningInsideDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if cgroup, err := os.ReadFile("/proc/1/cgroup"); err == nil && strings.Contains(string(cgroup), "docker") {
		return true
	}
	return false
}
