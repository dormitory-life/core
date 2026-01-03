package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/dormitory-life/core/internal/config"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
	"github.com/dormitory-life/utils/migrator"
)

type Driver interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// [Row]'s Scan method is called.
	// If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
	// Otherwise, [*Row.Scan] scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Database struct {
	db *sql.DB
}

type Repository interface {
	GetDormitories(ctx context.Context, request *dbtypes.GetDormitoriesRequest) (*dbtypes.GetDormitoriesResponse, error)
	GetDormitoryById(ctx context.Context, request *dbtypes.GetDormitoryByIdRequest) (*dbtypes.GetDormitoryByIdResponse, error)
	CreateDormitory(ctx context.Context, request *dbtypes.CreateDormitoryRequest) (*dbtypes.CreateDormitoryResponse, error)
	UpdateDormitory(ctx context.Context, request *dbtypes.UpdateDormitoryRequest) (*dbtypes.UpdateDormitoryResponse, error)
}

func New(db *sql.DB) Repository {
	return &Database{
		db: db,
	}
}

func InitDb(cfg config.DataBaseConfig) (*sql.DB, error) {
	connStr := cfg.GetConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if err := migrator.MigrateDB(connStr, cfg.MigrationsPath); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations completed successfully!")
	log.Println("Core service is ready")

	return db, nil
}
