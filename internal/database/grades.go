package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/dormitory-life/core/internal/constants"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

func (c *Database) GetDormitoriesAvgGrades(
	ctx context.Context,
	request *dbtypes.GetDormitoriesAvgGradesRequest,
) (*dbtypes.GetDormitoriesAvgGradesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", dberrors.ErrBadRequest)
	}

	resp, err := c.getDormitoriesAvgGrades(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getDormitoriesAvgGrades(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetDormitoriesAvgGradesRequest,
) (*dbtypes.GetDormitoriesAvgGradesResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		dormitoryAvgGradesTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryAvgGradesTableName)
	)

	queryBuilder := psql.Select("*").
		From(dormitoryAvgGradesTable).
		Where("period_date = (SELECT MAX(period_date) FROM dormitory_average_grades dag2 WHERE dag2.dormitory_id = dormitory_average_grades.dormitory_id)")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitories avg grades query: %v", dberrors.ErrInternal, err)
	}

	rows, err := driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing get dormitories avg grades query: %v", dberrors.ErrInternal, err)
	}
	defer rows.Close()

	var averageGrades []dbtypes.AvgGrade
	for rows.Next() {
		var avgGrade dbtypes.AvgGrade
		if err := rows.Scan(
			&avgGrade.Id,
			&avgGrade.DormitoryId,
			&avgGrade.PeriodDate,
			&avgGrade.AvgBathroomCleanliness,
			&avgGrade.AvgCorridorCleanliness,
			&avgGrade.AvgKitchenCleanliness,
			&avgGrade.AvgCleaningFrequency,
			&avgGrade.AvgRoomSpaciousness,
			&avgGrade.AvgCorridorSpaciousness,
			&avgGrade.AvgKitchenSpaciousness,
			&avgGrade.AvgShowerLocationConvenience,
			&avgGrade.AvgEquipmentMaintenance,
			&avgGrade.AvgWindowCondition,
			&avgGrade.AvgNoiseIsolation,
			&avgGrade.AvgCommonAreasEquipment,
			&avgGrade.AvgTransportAccessibility,
			&avgGrade.AvgAdministrationQuality,
			&avgGrade.AvgResidentsCultureLevel,
			&avgGrade.OverallAverage,
			&avgGrade.TotalRatings,
			&avgGrade.CreatedAt,
			&avgGrade.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: error scanning row: %v", dberrors.ErrInternal, err)
		}
		averageGrades = append(averageGrades, avgGrade)
	}

	return &dbtypes.GetDormitoriesAvgGradesResponse{
		Grades: averageGrades,
	}, nil
}

func (c *Database) GetDormitoryAvgGrades(
	ctx context.Context,
	request *dbtypes.GetDormitoryAvgGradesRequest,
) (*dbtypes.GetDormitoryAvgGradesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", dberrors.ErrBadRequest)
	}

	resp, err := c.getDormitoryAvgGrades(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) getDormitoryAvgGrades(
	ctx context.Context,
	driver Driver,
	request *dbtypes.GetDormitoryAvgGradesRequest,
) (*dbtypes.GetDormitoryAvgGradesResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		dormitoryAvgGradesTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.DormitoryAvgGradesTableName)
	)

	queryBuilder := psql.Select("*").
		From(dormitoryAvgGradesTable).
		Where(squirrel.Eq{"dormitory_id": request.DormitoryId}).
		OrderBy("period_date DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building get dormitories avg grades query: %v", dberrors.ErrInternal, err)
	}

	rows, err := driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: error executing get dormitories avg grades query: %v", dberrors.ErrInternal, err)
	}
	defer rows.Close()

	var averageGrades []dbtypes.AvgGrade
	for rows.Next() {
		var avgGrade dbtypes.AvgGrade
		if err := rows.Scan(
			&avgGrade.Id,
			&avgGrade.DormitoryId,
			&avgGrade.PeriodDate,
			&avgGrade.AvgBathroomCleanliness,
			&avgGrade.AvgCorridorCleanliness,
			&avgGrade.AvgKitchenCleanliness,
			&avgGrade.AvgCleaningFrequency,
			&avgGrade.AvgRoomSpaciousness,
			&avgGrade.AvgCorridorSpaciousness,
			&avgGrade.AvgKitchenSpaciousness,
			&avgGrade.AvgShowerLocationConvenience,
			&avgGrade.AvgEquipmentMaintenance,
			&avgGrade.AvgWindowCondition,
			&avgGrade.AvgNoiseIsolation,
			&avgGrade.AvgCommonAreasEquipment,
			&avgGrade.AvgTransportAccessibility,
			&avgGrade.AvgAdministrationQuality,
			&avgGrade.AvgResidentsCultureLevel,
			&avgGrade.OverallAverage,
			&avgGrade.TotalRatings,
			&avgGrade.CreatedAt,
			&avgGrade.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: error scanning row: %v", dberrors.ErrInternal, err)
		}
		averageGrades = append(averageGrades, avgGrade)
	}

	return &dbtypes.GetDormitoryAvgGradesResponse{
		Grades: averageGrades,
	}, nil
}

func (c *Database) CreateDormitoryGrade(
	ctx context.Context,
	request *dbtypes.CreateDormitoryGradeRequest,
) (*dbtypes.CreateDormitoryGradeResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", dberrors.ErrBadRequest)
	}

	resp, err := c.createDormitoryGrade(ctx, c.db, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Database) createDormitoryGrade(
	ctx context.Context,
	driver Driver,
	request *dbtypes.CreateDormitoryGradeRequest,
) (*dbtypes.CreateDormitoryGradeResponse, error) {
	if request == nil {
		return nil, dberrors.ErrBadRequest
	}

	var (
		psql                 = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
		dormitoryGradesTable = fmt.Sprintf("%s.%s", constants.SchemaName, constants.GradesTable)
	)

	queryBuilder := psql.Insert(dormitoryGradesTable).
		Columns(
			"dormitory_id",
			"user_id",
			"bathroom_cleanliness",
			"corridor_cleanliness",
			"kitchen_cleanliness",
			"cleaning_frequency",
			"room_spaciousness",
			"corridor_spaciousness",
			"kitchen_spaciousness",
			"shower_location_convenience",
			"equipment_maintenance",
			"window_condition",
			"noise_isolation",
			"common_areas_equipment",
			"transport_accessibility",
			"administration_quality",
			"residents_culture_level",
		).
		Values(
			request.DormitoryId,
			request.UserId,
			request.BathroomCleanliness,
			request.CorridorCleanliness,
			request.KitchenCleanliness,
			request.CleaningFrequency,
			request.RoomSpaciousness,
			request.CorridorSpaciousness,
			request.KitchenSpaciousness,
			request.ShowerLocationConvenience,
			request.EquipmentMaintenance,
			request.WindowCondition,
			request.NoiseIsolation,
			request.CommonAreasEquipment,
			request.TransportAccessibility,
			request.AdministrationQuality,
			request.ResidentsCultureLevel,
		).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: error building create grade query: %v", dberrors.ErrInternal, err)
	}

	row := driver.QueryRowContext(ctx, query, args...)

	var resp dbtypes.CreateDormitoryGradeResponse

	err = row.Scan(&resp.GradeId)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("%w: user can only rate this dormitory once per month", dberrors.ErrConflict)
		}

		return nil, fmt.Errorf("%w: error creating grade: %v", dberrors.ErrInternal, err)
	}

	return &resp, nil
}
