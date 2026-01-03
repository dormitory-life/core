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
			&dormitory.SupportEmail,
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
		&dormitory.SupportEmail,
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

func (c *Database) CreateDormitory(
	ctx context.Context,
	request *dbtypes.CreateDormitoryRequest,
) (*dbtypes.CreateDormitoryResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.createDormitory(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) createDormitory(
	ctx context.Context,
	driver Driver,
	request *dbtypes.CreateDormitoryRequest,
) (*dbtypes.CreateDormitoryResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		dormitoryTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryTableName)
	)

	queryBuilder := psql.Insert(dormitoryTable).
		Columns(
			"id", "name", "address", "support_email", "description",
		).
		Values(
			request.DormitoryId,
			request.Name,
			request.Address,
			request.SupportEmail,
			request.Description,
		).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building create dormitory query: %v", dberrors.ErrInternal, err)
	}

	row := driver.QueryRowContext(ctx, query, args...)

	var resp dbtypes.CreateDormitoryResponse
	err = row.Scan(
		&resp.DormitoryId,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: error scanning created dormitory: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}

func (c *Database) UpdateDormitory(
	ctx context.Context,
	request *dbtypes.UpdateDormitoryRequest,
) (*dbtypes.UpdateDormitoryResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.updateDormitory(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) updateDormitory(
	ctx context.Context,
	driver Driver,
	request *dbtypes.UpdateDormitoryRequest,
) (*dbtypes.UpdateDormitoryResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		dormitoryTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryTableName)
	)

	queryBuilder := psql.Update(dormitoryTable).Where(squirrel.Eq{"id": request.DormitoryId})

	setupUpdateFields(&queryBuilder, request)

	queryBuilder = queryBuilder.Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building update dormitory query: %v", dberrors.ErrInternal, err)
	}

	row := driver.QueryRowContext(ctx, query, args...)

	var resp dbtypes.UpdateDormitoryResponse
	err = row.Scan(
		&resp.DormitoryId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: dormitory not found", dberrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: error scanning updated dormitory: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}

func setupUpdateFields(
	queryBuilder *squirrel.UpdateBuilder,
	request *dbtypes.UpdateDormitoryRequest,
) {
	if request.Address != nil {
		*queryBuilder = queryBuilder.Set("address", request.Address)
	}

	if request.Name != nil {
		*queryBuilder = queryBuilder.Set("name", request.Name)
	}

	if request.SupportEmail != nil {
		*queryBuilder = queryBuilder.Set("support_email", request.SupportEmail)
	}

	if request.Description != nil {
		*queryBuilder = queryBuilder.Set("description", request.Description)
	}
}
