package database

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dormitory-life/core/internal/constants"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
	"github.com/google/uuid"
)

func (c *Database) GetChatMessages(
	ctx context.Context,
	request *dbtypes.GetChatMessagesRequest,
) (*dbtypes.GetChatMessagesResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.getChatMessages(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getChatMessages(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetChatMessagesRequest,
) (*dbtypes.GetChatMessagesResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		chatTable = fmt.Sprintf("%s.%s c", constants.SchemaName, constants.ChatTableName)
		userTable = fmt.Sprintf("%s.%s u", constants.SchemaName, constants.UsersTableName)
	)

	queryBuilder := psql.
		Select("c.id", "c.user_id", "c.dormitory_id", "c.text", "c.created_at", "u.email").
		From(chatTable).
		Join(fmt.Sprintf("%s ON u.id = c.user_id", userTable)).
		Where(squirrel.Eq{"c.dormitory_id": request.DormitoryID}).
		Limit(constants.DefaultPaginationPageSize).
		Offset(countOffset(request.Page, constants.DefaultPaginationPageSize))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get chat query: %v", dberrors.ErrInternal, err)
	}

	var messages []dbtypes.ChatMessage

	rows, err := driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing get chat query: %v", dberrors.ErrInternal, err)
	}

	defer rows.Close()

	for rows.Next() {
		var message dbtypes.ChatMessage

		if err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.DormitoryID,
			&message.Text,
			&message.CreatedAt,
			&message.Email,
		); err != nil {
			return nil, fmt.Errorf("%w: error scanning row: %v", dberrors.ErrInternal, err)
		}

		messages = append(messages, message)
	}

	return &dbtypes.GetChatMessagesResponse{
		Messages: messages,
	}, nil
}

func (c *Database) CreateChatMessage(
	ctx context.Context,
	request *dbtypes.CreateChatMessageRequest,
) (*dbtypes.CreateChatMessageResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.createChatMessage(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) createChatMessage(
	ctx context.Context,
	driver Driver,
	request *dbtypes.CreateChatMessageRequest,
) (*dbtypes.CreateChatMessageResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		chatTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.ChatTableName)
	)

	queryBuilder := psql.Insert(chatTable).
		Columns(
			"id", "user_id", "dormitory_id", "text", "created_at",
		).
		Values(
			uuid.New(),
			request.UserID,
			request.DormitoryID,
			request.Text,
			squirrel.Expr("now()"),
		).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building create message query: %v", dberrors.ErrInternal, err)
	}

	row := driver.QueryRowContext(ctx, query, args...)

	var resp dbtypes.CreateChatMessageResponse
	err = row.Scan(
		&resp.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: error scanning created message: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}
