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

func (c *Database) GetEmailsForSupport(
	ctx context.Context,
	request *dbtypes.GetEmailsForSupportRequest,
) (*dbtypes.GetEmailsForSupportResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		userTable      = fmt.Sprintf("%s.%s", constants.SchemaName, constants.UsersTableName)
		dormitoryTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryTableName)
	)

	queryBuilder := psql.
		Select("u.email", "d.support_email").
		From(fmt.Sprintf("%s u", userTable)).
		Join(fmt.Sprintf("%s d ON u.dormitory_id = d.id", dormitoryTable)).
		Where(squirrel.Eq{"u.id": request.UserId})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitories query: %v", dberrors.ErrInternal, err)
	}

	var userEmail, supportEmail string

	err = c.db.QueryRowContext(ctx, query, args...).Scan(&userEmail, &supportEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dberrors.ErrNotFound
		}

		return nil, fmt.Errorf("%w: error executing query: %v", dberrors.ErrInternal, err)
	}

	if supportEmail == "" {
		return nil, fmt.Errorf("%w: support email not configured for dormitory", dberrors.ErrNotFound)
	}

	return &dbtypes.GetEmailsForSupportResponse{
		UserEmail:    userEmail,
		SupportEmail: supportEmail,
	}, nil
}
