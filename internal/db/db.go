package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type DB struct {
	pool *pgxpool.Pool
}

var (
	instance    *DB
	once        sync.Once
	instanceErr error
)

func NewDB(dbFlags *config.DBFlags) (*DB, error) {
	once.Do(func() {
		connStr := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbFlags.Host, dbFlags.PORT, dbFlags.User, dbFlags.Password, dbFlags.Name, dbFlags.SSLMode,
		)

		poolConfig, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			instance = nil
			instanceErr = err
			return
		}

		poolConfig.MaxConns = 25
		poolConfig.MinConns = 5
		poolConfig.MaxConnLifetime = time.Hour
		poolConfig.MaxConnIdleTime = 30 * time.Minute
		poolConfig.HealthCheckPeriod = time.Minute

		pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			instance = nil
			instanceErr = err
			return
		}

		instance = &DB{pool}
		instanceErr = nil
	})

	return instance, instanceErr
}

func (db *DB) RegisterRepositories() *repository.Queries {
	return repository.New(db.pool)
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
