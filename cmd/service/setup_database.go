package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"narde/internal/config"
	"time"
)

func setupPostgresConnection(ctx context.Context, config config.PostgresCreds) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("can't open and connect to database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime * time.Second)
	db.SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("can't ping database: %w", err)
	}

	return db, nil
}
