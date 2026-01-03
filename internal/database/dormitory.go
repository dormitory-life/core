package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dormitory-life/core/internal/constants"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

func (c *Database) GetDormitories(
	ctx context.Context,
	request *dbtypes.GetDormitoriesRequest,
) (*dbtypes.GetDormitoriesResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.getDormitories(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getDormitories(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetDormitoriesRequest,
) (*dbtypes.GetDormitoriesResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		dormitoryTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryTableName)
	)

	queryBuilder := psql.
		Select("id", "name", "address", "support_email", "description").
		From(dormitoryTable)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitories query: %v", dberrors.ErrInternal, err)
	}

	var dormitories []dbtypes.Dormitory

	rows, err := driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing get dormitories query: %v", dberrors.ErrInternal, err)
	}

	defer rows.Close()

	for rows.Next() {
		var dormitory dbtypes.Dormitory

		if err := rows.Scan(
			&dormitory.Id,
			&dormitory.Name,
			&dormitory.Address,
			&dormitory.Support_email,
			&dormitory.Description,
		); err != nil {
			return nil, fmt.Errorf("%w: error scanning row: %v", dberrors.ErrInternal, err)
		}

		dormitories = append(dormitories, dormitory)
	}

	return &dbtypes.GetDormitoriesResponse{
		Dormitories: dormitories,
	}, nil
}

func (c *Database) GetDormitoryById(
	ctx context.Context,
	request *dbtypes.GetDormitoryByIdRequest,
) (*dbtypes.GetDormitoryByIdResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.getDormitoryById(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getDormitoryById(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetDormitoryByIdRequest,
) (*dbtypes.GetDormitoryByIdResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		dormitoryTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryTableName)
	)

	queryBuilder := psql.
		Select("id", "name", "address", "support_email", "description").
		From(dormitoryTable).
		Where(squirrel.Eq{"id": request.DormitoryId}).
		Limit(1)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitory by id query: %v", dberrors.ErrInternal, err)
	}

	var dormitory dbtypes.Dormitory

	err = driver.QueryRowContext(ctx, query, args...).Scan(
		&dormitory.Id,
		&dormitory.Name,
		&dormitory.Address,
		&dormitory.Support_email,
		&dormitory.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: dormitory not found", dberrors.ErrNotFound)
		}

		return nil, fmt.Errorf("%w: error executing get dormitory by id query: %v", dberrors.ErrInternal, err)
	}

	return &dbtypes.GetDormitoryByIdResponse{
		Dormitory: dormitory,
	}, nil
}
