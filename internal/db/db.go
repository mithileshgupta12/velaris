package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/mithileshgupta12/velaris/internal/config"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbFlags *config.DBFlags) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbFlags.Host, dbFlags.PORT, dbFlags.User, dbFlags.Password, dbFlags.Name, dbFlags.SSLMode,
	)

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(15 * time.Minute)
	conn.SetConnMaxIdleTime(5 * time.Minute)

	return &DB{
		conn,
	}, nil
}

func (db *DB) RegisterRepositories() {
	//
}

func (db *DB) Ping() error {
	if err := db.conn.Ping(); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
