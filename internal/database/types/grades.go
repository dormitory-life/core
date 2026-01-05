package types

import "time"

type Grade struct {
	Id          string
	DormitoryId string
	UserId      string

	BathroomCleanliness       int
	CorridorCleanliness       int
	KitchenCleanliness        int
	CleaningFrequency         int
	RoomSpaciousness          int
	CorridorSpaciousness      int
	KitchenSpaciousness       int
	ShowerLocationConvenience int
	EquipmentMaintenance      int
	WindowCondition           int
	NoiseIsolation            int
	CommonAreasEquipment      int
	TransportAccessibility    int
	AdministrationQuality     int
	ResidentsCultureLevel     int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type (
	CreateDormitoryGradeRequest struct {
		DormitoryId string
		UserId      string

		BathroomCleanliness       int
		CorridorCleanliness       int
		KitchenCleanliness        int
		CleaningFrequency         int
		RoomSpaciousness          int
		CorridorSpaciousness      int
		KitchenSpaciousness       int
		ShowerLocationConvenience int
		EquipmentMaintenance      int
		WindowCondition           int
		NoiseIsolation            int
		CommonAreasEquipment      int
		TransportAccessibility    int
		AdministrationQuality     int
		ResidentsCultureLevel     int
	}

	CreateDormitoryGradeResponse struct {
		GradeId string
	}
)

type AvgGrade struct {
	Id                           string
	DormitoryId                  string
	PeriodDate                   time.Time
	AvgBathroomCleanliness       float64
	AvgCorridorCleanliness       float64
	AvgKitchenCleanliness        float64
	AvgCleaningFrequency         float64
	AvgRoomSpaciousness          float64
	AvgCorridorSpaciousness      float64
	AvgKitchenSpaciousness       float64
	AvgShowerLocationConvenience float64
	AvgEquipmentMaintenance      float64
	AvgWindowCondition           float64
	AvgNoiseIsolation            float64
	AvgCommonAreasEquipment      float64
	AvgTransportAccessibility    float64
	AvgAdministrationQuality     float64
	AvgResidentsCultureLevel     float64
	OverallAverage               float64
	TotalRatings                 int
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

type (
	GetDormitoriesAvgGradesRequest  struct{}
	GetDormitoriesAvgGradesResponse struct {
		Grades []AvgGrade
	}
)

type (
	GetDormitoryAvgGradesRequest struct {
		DormitoryId string
	}
	GetDormitoryAvgGradesResponse struct {
		Grades []AvgGrade
	}
)
