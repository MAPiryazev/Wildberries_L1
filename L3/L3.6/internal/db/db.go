package db

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wb-go/wbf/dbpg"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.6/internal/config"
)

func Connect(cfg *config.Config) (*dbpg.DB, error) {
	masterDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	opts := &dbpg.Options{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := dbpg.New(masterDSN, nil, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Master.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *dbpg.DB) error {
	dir := "migrations"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if filepath.Ext(name) != ".sql" {
			continue
		}

		files = append(files, filepath.Join(dir, name))
	}

	sort.Strings(files)

	for _, path := range files {
		if err := runMigrationFile(db, path); err != nil {
			return err
		}
	}

	return nil
}

func runMigrationFile(db *dbpg.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("migration file not found %s: %w", path, err)
		}
		return fmt.Errorf("read migration file %s: %w", path, err)
	}

	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := db.ExecContext(context.Background(), stmt); err != nil {
			return fmt.Errorf("migration %s error: %w", path, err)
		}
	}

	return nil
}
