package database

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dormitory-life/core/internal/constants"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

func (c *Database) GetDormitoryEvents(
	ctx context.Context,
	request *dbtypes.GetDormitoryEventsRequest,
) (*dbtypes.GetDormitoryEventsResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.getDormitoryEvents(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getDormitoryEvents(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetDormitoryEventsRequest,
) (*dbtypes.GetDormitoryEventsResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		feedTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.FeedTableName)
	)

	queryBuilder := psql.
		Select(
			"id", "dormitory_id", "title", "description", "created_at",
		).
		From(feedTable).
		Where(squirrel.Eq{"dormitory_id": request.DormitoryId}).
		Offset(countOffset(request.Page, constants.DefaultEventsPageSize)).
		Limit(constants.DefaultEventsPageSize).
		OrderBy("created_at DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitory events query: %v", dberrors.ErrInternal, err)
	}

	var events []dbtypes.Event

	rows, err := driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing get dormitory events query: %v", dberrors.ErrInternal, err)
	}

	defer rows.Close()

	for rows.Next() {
		var event dbtypes.Event

		if err := rows.Scan(
			&event.EventId,
			&event.DormitoryId,
			&event.Title,
			&event.Description,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: error scanning row: %v", dberrors.ErrInternal, err)
		}

		events = append(events, event)
	}

	return &dbtypes.GetDormitoryEventsResponse{
		Events: events,
	}, nil
}

func (c *Database) CreateDormitoryEvent(
	ctx context.Context,
	request *dbtypes.CreateDormitoryEventRequest,
) (*dbtypes.CreateDormitoryEventResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.createDormitoryEvent(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) createDormitoryEvent(
	ctx context.Context,
	driver Driver,
	request *dbtypes.CreateDormitoryEventRequest,
) (*dbtypes.CreateDormitoryEventResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		feedTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.FeedTableName)
	)

	queryBuilder := psql.Insert(feedTable).
		Columns(
			"dormitory_id", "title", "description",
		).
		Values(
			request.DormitoryId,
			request.Title,
			request.Description,
		).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building create event query: %v", dberrors.ErrInternal, err)
	}

	row := driver.QueryRowContext(ctx, query, args...)

	var resp dbtypes.CreateDormitoryEventResponse
	err = row.Scan(
		&resp.EventId,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: error scanning created event: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}

func (c *Database) DeleteDormitoryEvent(
	ctx context.Context,
	request *dbtypes.DeleteDormitoryEventRequest,
) (*dbtypes.DeleteDormitoryEventResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	resp, err := c.deleteDormitoryEvent(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) deleteDormitoryEvent(
	ctx context.Context,
	driver Driver,
	request *dbtypes.DeleteDormitoryEventRequest,
) (*dbtypes.DeleteDormitoryEventResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		feedTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.FeedTableName)
	)

	queryBuilder := psql.Delete(feedTable).
		Where(squirrel.Eq{"id": request.EventId})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building delete event query: %v", dberrors.ErrInternal, err)
	}

	_, err = driver.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing delete event query: %v", dberrors.ErrInternal, err)
	}

	return &dbtypes.DeleteDormitoryEventResponse{}, nil
}
