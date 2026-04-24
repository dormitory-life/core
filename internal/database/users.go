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

func (c *Database) GetUsersRole(
	ctx context.Context,
	request *dbtypes.GetUsersRoleRequest,
) (*dbtypes.GetUsersRoleResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.getUsersRole(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getUsersRole(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetUsersRoleRequest,
) (*dbtypes.GetUsersRoleResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		userTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.UsersTableName)
	)

	queryBuilder := psql.
		Select("role").
		From(userTable).
		Where(squirrel.Eq{"id": request.UserId}).
		Limit(1)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get users role query: %v", dberrors.ErrInternal, err)
	}

	var resp dbtypes.GetUsersRoleResponse

	err = driver.QueryRowContext(ctx, query, args...).Scan(
		&resp.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: user not found", dberrors.ErrNotFound)
		}

		return nil, fmt.Errorf("%w: error executing get users role query: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}
