package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(dbFlags *config.DBFlags) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbFlags.Host, dbFlags.PORT, dbFlags.User, dbFlags.Password, dbFlags.Name, dbFlags.SSLMode,
	)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn,
	}, nil
}

func (db *DB) RegisterRepositories() *repository.Queries {
	return repository.New(db.conn)
}

func (db *DB) GetConn() *pgx.Conn {
	return db.conn
}

func (db *DB) Ping() error {
	if err := db.conn.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close(context.Background())
}
