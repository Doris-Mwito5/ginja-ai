package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Doris-Mwito5/ginja-ai/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	DEFAULT_MAX_CONNECTIONS  = 40
	MAX_CONNECTION_IDLE_TIME = 5 * time.Minute
)

var db DB

type SQLOperations interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ValidForPostgres() bool
}

type pgSQLOperations struct {
	*sql.Tx
}

func (o *pgSQLOperations) ValidForPostgres() bool {
	return true
}

type DB interface {
	SQLOperations
	Begin() (*sql.Tx, error)
	Close() error
	Ping() error
	InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error
	Valid() bool
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type appDB struct {
	*sql.DB
	pool  *pgxpool.Pool
	valid bool
}

func (db *appDB) InTransaction(ctx context.Context, operations func(context.Context, SQLOperations) error) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sqlOperations := &pgSQLOperations{
		Tx: tx,
	}

	if err = operations(ctx, sqlOperations); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
}

func (db *appDB) ValidForPostgres() bool {
	return true
}

func (db *appDB) Valid() bool {
	return db.valid
}

func InitDB(databaseURL string) DB {
	return InitDBWithURL(
		databaseURL,
	)
}

func InitDBWithURL(databaseURL string) DB {

	pgDB, pool := newPostgresDBWithURL(databaseURL)

	db = &appDB{
		DB:    pgDB,
		pool:  pool,
		valid: true,
	}

	err := db.Ping()
	if err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	return db
}

func GetDB() DB {
	return db
}

func newPostgresDBWithURL(databaseURL string) (*sql.DB, *pgxpool.Pool) {
	if databaseURL == "" {
		logger.Fatal("database url is empty and is required")
	}

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		logger.Fatalf("unable to create connection pool: [%v]", err)
	}

	// Configure connection pool
	pool.Config().MaxConns = DEFAULT_MAX_CONNECTIONS
	pool.Config().MaxConnIdleTime = MAX_CONNECTION_IDLE_TIME

	// Verify the pool connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatalf("Database pool connection failed: [%v]", err)
	}

	// Also create sql.DB for backward compatibility
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		logger.Fatalf("sql.Open function call failed: [%v]", err)
	}

	// Set connection pool configurations for sql.DB
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	// Verify the database connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Database connection failed: [%v]", err)
	}

	return db, pool
}
