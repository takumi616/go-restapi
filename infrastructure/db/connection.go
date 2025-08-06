package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/takumi616/go-restapi/shared/config"
)

func NewDBConnection(ctx context.Context, dbConf *config.DatabaseConfig) (*sql.DB, error) {
	// Create datasourcename, using database connection config
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConf.Connection.Host,
		dbConf.Connection.Port,
		dbConf.Connection.User,
		dbConf.Connection.Password,
		dbConf.Connection.DbName,
		dbConf.Connection.Sslmode,
	)

	db, err := sql.Open(dbConf.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	// Set connection pool config
	db.SetMaxOpenConns(dbConf.Pool.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.Pool.MaxIdleConns)
	db.SetConnMaxLifetime(dbConf.Pool.ConnMaxLifetime)
	db.SetConnMaxIdleTime(dbConf.Pool.ConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if database connection is alive
	err = db.PingContext(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("Database connection is not alive: %w", err)
	}

	return db, nil
}
