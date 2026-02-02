package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	cfg "shortener/internal/config"
	httpapi "shortener/internal/handlers"
	"shortener/internal/repository/clickhouse"
	"shortener/internal/repository/psql"
	"shortener/internal/service"
	"shortener/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// API config
	apiCfg := cfg.LoadAPIConfig("../../environment/.env")

	// PSQL config and pool
	dbCfg, err := cfg.LoadDBPSQLConfig("../../environment/.env")
	if err != nil {
		log.Fatalf("db config: %v", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pool.Close()
	urlRepo := psql.NewRepo(pool)

	// ClickHouse config and repo (HTTP)
	chCfg, err := cfg.LoadClickHouseConfig("../../environment/.env")
	if err != nil {
		log.Fatalf("clickhouse config: %v", err)
	}
	chRepo := clickhouse.NewAnalyticsRepo(chCfg.Host, chCfg.HTTPPort, chCfg.Database, chCfg.User, chCfg.Password, nil)

	// Service
	svc := service.NewService(urlRepo, chRepo, func() (string, error) { return utils.GenerateShortCode(8) })

	// HTTP
	h := httpapi.NewHandler(svc)
	addr := ":" + apiCfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, h.Router()); err != nil {
		log.Fatal(err)
	}
}
