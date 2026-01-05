package requestmodels

import (
	"time"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

type Grade struct {
	Id          string `json:"id"`
	DormitoryId string `json:"dormitory_id"`
	UserId      string `json:"user_id"`

	BathroomCleanliness       int `json:"bathroom_cleanliness"`
	CorridorCleanliness       int `json:"corridor_cleanliness"`
	KitchenCleanliness        int `json:"kitchen_cleanliness"`
	CleaningFrequency         int `json:"cleaning_frequency"`
	RoomSpaciousness          int `json:"room_spaciousness"`
	CorridorSpaciousness      int `json:"corridor_spaciousness"`
	KitchenSpaciousness       int `json:"kitchen_spaciousness"`
	ShowerLocationConvenience int `json:"shower_location_convenience"`
	EquipmentMaintenance      int `json:"equipment_maintenance"`
	WindowCondition           int `json:"window_condition"`
	NoiseIsolation            int `json:"noise_isolation"`
	CommonAreasEquipment      int `json:"common_areas_equipment"`
	TransportAccessibility    int `json:"transport_accessibility"`
	AdministrationQuality     int `json:"administration_quality"`
	ResidentsCultureLevel     int `json:"residents_culture_level"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	CreateDormitoryGradeRequest struct {
		DormitoryId string

		BathroomCleanliness       int `json:"bathroom_cleanliness"`
		CorridorCleanliness       int `json:"corridor_cleanliness"`
		KitchenCleanliness        int `json:"kitchen_cleanliness"`
		CleaningFrequency         int `json:"cleaning_frequency"`
		RoomSpaciousness          int `json:"room_spaciousness"`
		CorridorSpaciousness      int `json:"corridor_spaciousness"`
		KitchenSpaciousness       int `json:"kitchen_spaciousness"`
		ShowerLocationConvenience int `json:"shower_location_convenience"`
		EquipmentMaintenance      int `json:"equipment_maintenance"`
		WindowCondition           int `json:"window_condition"`
		NoiseIsolation            int `json:"noise_isolation"`
		CommonAreasEquipment      int `json:"common_areas_equipment"`
		TransportAccessibility    int `json:"transport_accessibility"`
		AdministrationQuality     int `json:"administration_quality"`
		ResidentsCultureLevel     int `json:"residents_culture_level"`
	}

	CreateDormitoryGradeResponse struct {
		GradeId string `json:"grade_id"`
	}
)

func (r *CreateDormitoryGradeResponse) From(msg *dbtypes.CreateDormitoryGradeResponse) *CreateDormitoryGradeResponse {
	if msg == nil {
		return nil
	}

	res := &CreateDormitoryGradeResponse{
		GradeId: msg.GradeId,
	}

	return res
}

type AvgGrade struct {
	Id          string    `json:"id"`
	DormitoryId string    `json:"dormitory_id"`
	PeriodDate  time.Time `json:"period_date"`

	AvgBathroomCleanliness       float64 `json:"avg_bathroom_cleanliness"`
	AvgCorridorCleanliness       float64 `json:"avg_corridor_cleanliness"`
	AvgKitchenCleanliness        float64 `json:"avg_kitchen_cleanliness"`
	AvgCleaningFrequency         float64 `json:"avg_cleaning_frequency"`
	AvgRoomSpaciousness          float64 `json:"avg_room_spaciousness"`
	AvgCorridorSpaciousness      float64 `json:"avg_corridor_spaciousness"`
	AvgKitchenSpaciousness       float64 `json:"avg_kitchen_spaciousness"`
	AvgShowerLocationConvenience float64 `json:"avg_shower_location_convenience"`
	AvgEquipmentMaintenance      float64 `json:"avg_equipment_maintenance"`
	AvgWindowCondition           float64 `json:"avg_window_condition"`
	AvgNoiseIsolation            float64 `json:"avg_noise_isolation"`
	AvgCommonAreasEquipment      float64 `json:"avg_common_areas_equipment"`
	AvgTransportAccessibility    float64 `json:"avg_transport_accessibility"`
	AvgAdministrationQuality     float64 `json:"avg_administration_quality"`
	AvgResidentsCultureLevel     float64 `json:"avg_residents_culture_level"`

	OverallAverage float64 `json:"overall_average"`
	TotalRatings   int     `json:"total_ratings"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type (
	GetDormitoriesAvgGradesRequest  struct{}
	GetDormitoriesAvgGradesResponse struct {
		Grades []AvgGrade `json:"avg_grades"`
	}
)

func (r *AvgGrade) From(msg *dbtypes.AvgGrade) *AvgGrade {
	if msg == nil {
		return nil
	}

	return &AvgGrade{
		Id:                           msg.Id,
		DormitoryId:                  msg.DormitoryId,
		PeriodDate:                   msg.PeriodDate,
		AvgBathroomCleanliness:       msg.AvgBathroomCleanliness,
		AvgCorridorCleanliness:       msg.AvgCorridorCleanliness,
		AvgKitchenCleanliness:        msg.AvgKitchenCleanliness,
		AvgCleaningFrequency:         msg.AvgCleaningFrequency,
		AvgRoomSpaciousness:          msg.AvgRoomSpaciousness,
		AvgCorridorSpaciousness:      msg.AvgCorridorSpaciousness,
		AvgKitchenSpaciousness:       msg.AvgKitchenSpaciousness,
		AvgShowerLocationConvenience: msg.AvgShowerLocationConvenience,
		AvgEquipmentMaintenance:      msg.AvgEquipmentMaintenance,
		AvgWindowCondition:           msg.AvgWindowCondition,
		AvgNoiseIsolation:            msg.AvgNoiseIsolation,
		AvgCommonAreasEquipment:      msg.AvgCommonAreasEquipment,
		AvgTransportAccessibility:    msg.AvgTransportAccessibility,
		AvgAdministrationQuality:     msg.AvgAdministrationQuality,
		AvgResidentsCultureLevel:     msg.AvgResidentsCultureLevel,
		OverallAverage:               msg.OverallAverage,
		TotalRatings:                 msg.TotalRatings,
		CreatedAt:                    msg.CreatedAt,
		UpdatedAt:                    msg.UpdatedAt,
	}
}

func (r *GetDormitoriesAvgGradesResponse) From(msg *dbtypes.GetDormitoriesAvgGradesResponse) *GetDormitoriesAvgGradesResponse {
	if msg == nil {
		return nil
	}

	res := &GetDormitoriesAvgGradesResponse{
		Grades: make([]AvgGrade, 0),
	}

	for _, val := range msg.Grades {
		avgGrade := *new(AvgGrade).From(&val)
		res.Grades = append(res.Grades, avgGrade)
	}

	return res
}

type (
	GetDormitoryAvgGradesRequest struct {
		DormitoryId string `json:"dormitory_id"`
	}
	GetDormitoryAvgGradesResponse struct {
		Grades []AvgGrade `json:"avg_grades"`
	}
)

func (r *GetDormitoryAvgGradesResponse) From(msg *dbtypes.GetDormitoryAvgGradesResponse) *GetDormitoryAvgGradesResponse {
	if msg == nil {
		return nil
	}

	res := &GetDormitoryAvgGradesResponse{
		Grades: make([]AvgGrade, 0),
	}

	for _, val := range msg.Grades {
		avgGrade := *new(AvgGrade).From(&val)
		res.Grades = append(res.Grades, avgGrade)
	}

	return res
}
