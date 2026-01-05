CREATE TABLE IF NOT EXISTS dormitory_average_grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    dormitory_id VARCHAR(2) NOT NULL REFERENCES dormitory (id) ON DELETE CASCADE,
    period_date DATE NOT NULL, -- Дата периода (первое число месяца: '2025-01-01')
    avg_bathroom_cleanliness DECIMAL(3, 2) NOT NULL,
    avg_corridor_cleanliness DECIMAL(3, 2) NOT NULL,
    avg_kitchen_cleanliness DECIMAL(3, 2) NOT NULL,
    avg_cleaning_frequency DECIMAL(3, 2) NOT NULL,
    avg_room_spaciousness DECIMAL(3, 2) NOT NULL,
    avg_corridor_spaciousness DECIMAL(3, 2) NOT NULL,
    avg_kitchen_spaciousness DECIMAL(3, 2) NOT NULL,
    avg_shower_location_convenience DECIMAL(3, 2) NOT NULL,
    avg_equipment_maintenance DECIMAL(3, 2) NOT NULL,
    avg_window_condition DECIMAL(3, 2) NOT NULL,
    avg_noise_isolation DECIMAL(3, 2) NOT NULL,
    avg_common_areas_equipment DECIMAL(3, 2) NOT NULL,
    avg_transport_accessibility DECIMAL(3, 2) NOT NULL,
    avg_administration_quality DECIMAL(3, 2) NOT NULL,
    avg_residents_culture_level DECIMAL(3, 2) NOT NULL,
    overall_average DECIMAL(3, 2) NOT NULL,
    total_ratings INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (dormitory_id, period_date)
);

CREATE INDEX IF NOT EXISTS idx_avg_grades_dormitory_id ON dormitory_average_grades (dormitory_id);

CREATE INDEX IF NOT EXISTS idx_avg_grades_period_date ON dormitory_average_grades (period_date DESC);

CREATE INDEX IF NOT EXISTS idx_avg_grades_dormitory_period ON dormitory_average_grades (
    dormitory_id,
    period_date DESC
);
