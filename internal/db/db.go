package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(dbFlags *config.DBFlags) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbFlags.Host, dbFlags.PORT, dbFlags.User, dbFlags.Password, dbFlags.Name, dbFlags.SSLMode,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &DB{
		pool,
	}, nil
}

func (db *DB) RegisterRepositories() *repository.Queries {
	return repository.New(db.pool)
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Ping() error {
	if err := db.pool.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() {
	db.pool.Close()
}
