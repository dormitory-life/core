package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/dormitory-life/core/internal/config"
	"github.com/dormitory-life/core/internal/constants"
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

	GetDormitoriesAvgGrades(ctx context.Context, request *dbtypes.GetDormitoriesAvgGradesRequest) (*dbtypes.GetDormitoriesAvgGradesResponse, error)
	GetDormitoryAvgGrades(ctx context.Context, request *dbtypes.GetDormitoryAvgGradesRequest) (*dbtypes.GetDormitoryAvgGradesResponse, error)
	CreateDormitoryGrade(ctx context.Context, request *dbtypes.CreateDormitoryGradeRequest) (*dbtypes.CreateDormitoryGradeResponse, error)

	GetEmailsForSupport(ctx context.Context, request *dbtypes.GetEmailsForSupportRequest) (*dbtypes.GetEmailsForSupportResponse, error)

	GetReviews(ctx context.Context, request *dbtypes.GetDormitoryReviewsRequest) (*dbtypes.GetDormitoryReviewsResponse, error)
	CreateReview(ctx context.Context, request *dbtypes.CreateReviewRequest) (*dbtypes.CreateReviewResponse, error)
	DeleteReview(ctx context.Context, request *dbtypes.DeleteReviewRequest) (*dbtypes.DeleteReviewResponse, error)
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

func countOffset(page, pageSize uint64) uint64 {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = constants.DefaultPaginationPageSize
	}

	return (page - 1) * pageSize
}
