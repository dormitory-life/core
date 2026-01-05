package core

import (
	"context"
	"fmt"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
)

func (s *CoreService) GetDormitoriesAvgGrades(
	ctx context.Context,
	request *rmodel.GetDormitoriesAvgGradesRequest,
) (*rmodel.GetDormitoriesAvgGradesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	resp, err := s.repository.GetDormitoriesAvgGrades(ctx, &dbtypes.GetDormitoriesAvgGradesRequest{})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting dormitories avg grades: %v", ErrInternal, err)
	}

	res := new(rmodel.GetDormitoriesAvgGradesResponse).From(resp)

	return res, nil
}

func (s *CoreService) GetDormitoryAvgGrades(
	ctx context.Context,
	request *rmodel.GetDormitoryAvgGradesRequest,
) (*rmodel.GetDormitoryAvgGradesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	resp, err := s.repository.GetDormitoryAvgGrades(ctx, &dbtypes.GetDormitoryAvgGradesRequest{
		DormitoryId: request.DormitoryId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting dormitories avg grades: %v", ErrInternal, err)
	}

	res := new(rmodel.GetDormitoryAvgGradesResponse).From(resp)

	return res, nil
}

func (s *CoreService) CreateDormitoryGrade(
	ctx context.Context,
	request *rmodel.CreateDormitoryGradeRequest,
) (*rmodel.CreateDormitoryGradeResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	userId, dormitoryId, err := s.extractIdsFromRequestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting ids from context: %v", ErrInternal, err)
	}

	if err := s.checkAccess(
		ctx,
		&rmodel.CheckAccessRequest{
			UserId:       userId,
			DormitoryId:  request.DormitoryId,
			RoleRequired: false,
		},
	); err != nil {
		return nil, err
	}

	resp, err := s.repository.CreateDormitoryGrade(ctx, &dbtypes.CreateDormitoryGradeRequest{
		DormitoryId:               dormitoryId,
		UserId:                    userId,
		BathroomCleanliness:       request.BathroomCleanliness,
		CorridorCleanliness:       request.CorridorCleanliness,
		KitchenCleanliness:        request.KitchenCleanliness,
		CleaningFrequency:         request.CleaningFrequency,
		RoomSpaciousness:          request.RoomSpaciousness,
		CorridorSpaciousness:      request.CorridorSpaciousness,
		KitchenSpaciousness:       request.KitchenSpaciousness,
		ShowerLocationConvenience: request.ShowerLocationConvenience,
		EquipmentMaintenance:      request.EquipmentMaintenance,
		WindowCondition:           request.WindowCondition,
		NoiseIsolation:            request.NoiseIsolation,
		CommonAreasEquipment:      request.CommonAreasEquipment,
		TransportAccessibility:    request.TransportAccessibility,
		AdministrationQuality:     request.AdministrationQuality,
		ResidentsCultureLevel:     request.ResidentsCultureLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error creating grade: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.CreateDormitoryGradeResponse).From(resp)

	return res, nil
}
