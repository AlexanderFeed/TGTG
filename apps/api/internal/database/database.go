package database

// Package database owns infrastructure-level PostgreSQL setup. SQL for actual
// features lives in package store; this package only opens the pool and brings
// the database schema up to date.

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eshche-est/eshche-est/apps/api/migrations"
	// The blank import runs pgx's init function, which registers the "pgx"
	// driver name used by sql.Open below. We do not call pgx directly here.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Open creates and verifies a PostgreSQL connection pool. sql.Open itself is
// lazy, so PingContext is necessary to prove that the URL and server work now.
func Open(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	// Pool limits prevent one API process from exhausting PostgreSQL's available
	// connections under load. These values are intentionally small for a pet app.
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return db, nil
}

// Migrate applies every not-yet-applied SQL migration in version order. Goose
// records completed versions in its own table, so restarting is idempotent.
func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations.Files)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
