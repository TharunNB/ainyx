package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, cfg *Config) *pgxpool.Pool {

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Unable to parse databse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database connection ping failed: %v", err)
	}

	fmt.Println("COnnection to PostgreSQL connected using pgxpool")
	return pool
}
