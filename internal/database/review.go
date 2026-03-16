package database

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dormitory-life/core/internal/constants"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

func (c *Database) GetReviews(
	ctx context.Context,
	request *dbtypes.GetDormitoryReviewsRequest,
) (*dbtypes.GetDormitoryReviewsResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.getReviews(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getReviews(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetDormitoryReviewsRequest,
) (*dbtypes.GetDormitoryReviewsResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		reviewTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.ReviewTableName)
	)

	queryBuilder := psql.
		Select(
			"id", "owner_id", "dormitory_id", "title", "description", "created_at",
		).
		From(reviewTable).
		Where(squirrel.Eq{"dormitory_id": request.DormitoryId}).
		Offset(countOffset(request.Page, constants.DefaultReviewsPageSize)).
		Limit(constants.DefaultReviewsPageSize).
		OrderBy("created_at DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitory reviews query: %v", dberrors.ErrInternal, err)
	}

	var reviews []dbtypes.Review

	rows, err := driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing get dormitory reviews query: %v", dberrors.ErrInternal, err)
	}

	defer rows.Close()

	for rows.Next() {
		var review dbtypes.Review

		if err := rows.Scan(
			&review.ReviewId,
			&review.OwnerId,
			&review.DormitoryId,
			&review.Title,
			&review.Description,
			&review.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: error scanning row: %v", dberrors.ErrInternal, err)
		}

		reviews = append(reviews, review)
	}

	return &dbtypes.GetDormitoryReviewsResponse{
		Reviews: reviews,
	}, nil
}

func (c *Database) CreateReview(
	ctx context.Context,
	request *dbtypes.CreateReviewRequest,
) (*dbtypes.CreateReviewResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.createReview(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) createReview(
	ctx context.Context,
	driver Driver,
	request *dbtypes.CreateReviewRequest,
) (*dbtypes.CreateReviewResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		reviewTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.ReviewTableName)
	)

	queryBuilder := psql.Insert(reviewTable).
		Columns(
			"owner_id", "dormitory_id", "title", "description",
		).
		Values(
			request.OwnerId,
			request.DormitoryId,
			request.Title,
			request.Description,
		).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building create review query: %v", dberrors.ErrInternal, err)
	}

	row := driver.QueryRowContext(ctx, query, args...)

	var resp dbtypes.CreateReviewResponse
	err = row.Scan(
		&resp.ReviewId,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: error scanning created review: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}

func (c *Database) DeleteReview(
	ctx context.Context,
	request *dbtypes.DeleteReviewRequest,
) (*dbtypes.DeleteReviewResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.deleteReview(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) deleteReview(
	ctx context.Context,
	driver Driver,
	request *dbtypes.DeleteReviewRequest,
) (*dbtypes.DeleteReviewResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		reviewTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.ReviewTableName)
	)

	queryBuilder := psql.Delete(reviewTable).
		Where(squirrel.Eq{"id": request.ReviewId})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building delete review query: %v", dberrors.ErrInternal, err)
	}

	_, err = driver.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing delete review query: %v", dberrors.ErrInternal, err)
	}

	return &dbtypes.DeleteReviewResponse{}, nil
}
