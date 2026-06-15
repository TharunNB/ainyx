package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	fmt.Println("Connection to PostgreSQL established using pgxpool")
	return pool
}

func RunMigrations(pool *pgxpool.Pool) error {
	migrationSQL, err := os.ReadFile("db/migrations/001_init_schema.sql")
	if err != nil {
		return fmt.Errorf("reading migration file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := pool.Exec(ctx, string(migrationSQL)); err != nil {
		return fmt.Errorf("executing migration: %w", err)
	}
	return nil
}
