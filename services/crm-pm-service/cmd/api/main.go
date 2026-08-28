package main

import (
	"log"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/cmd/config"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/cmd/server"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/db"
	"github.com/erickgreco/startup-backend/shared/logger"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.LogLevel)

	dbcfg := db.Config{
		MaxConns:              cfg.DBMaxConns,
		MinConns:              cfg.DBMinConns,
		MaxConnLifeTime:       cfg.DBMaxConnLifeTime,
		MaxConnIdleTime:       cfg.DBMaxConnIdleTime,
		HealthCheckPeriod:     cfg.DBHealthCheckPeriod,
		MaxConnLifeTimeJitter: cfg.DBMaxConnLifeTimeJitter,
	}

	dbpool, err := db.Connect(cfg.DBURL, dbcfg)
	if err != nil {
		log.Fatalf("db_conn_err: %v", err)
	}
	defer dbpool.Close()
	logger.Info("database_connection_established")

	srv := server.New(dbpool, cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("server_failed: %v", err)
	}
}
