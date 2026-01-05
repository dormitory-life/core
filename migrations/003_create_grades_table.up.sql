-- core/migrations/003_create_grades_table.sql
CREATE TABLE IF NOT EXISTS grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    dormitory_id VARCHAR(2) NOT NULL REFERENCES dormitory (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    -- Оценочные критерии (1-5)
    bathroom_cleanliness INTEGER NOT NULL CHECK (
        bathroom_cleanliness BETWEEN 1 AND 5
    ),
    corridor_cleanliness INTEGER NOT NULL CHECK (
        corridor_cleanliness BETWEEN 1 AND 5
    ),
    kitchen_cleanliness INTEGER NOT NULL CHECK (
        kitchen_cleanliness BETWEEN 1 AND 5
    ),
    cleaning_frequency INTEGER NOT NULL CHECK (
        cleaning_frequency BETWEEN 1 AND 5
    ),
    room_spaciousness INTEGER NOT NULL CHECK (
        room_spaciousness BETWEEN 1 AND 5
    ),
    corridor_spaciousness INTEGER NOT NULL CHECK (
        corridor_spaciousness BETWEEN 1 AND 5
    ),
    kitchen_spaciousness INTEGER NOT NULL CHECK (
        kitchen_spaciousness BETWEEN 1 AND 5
    ),
    shower_location_convenience INTEGER NOT NULL CHECK (
        shower_location_convenience BETWEEN 1 AND 5
    ),
    equipment_maintenance INTEGER NOT NULL CHECK (
        equipment_maintenance BETWEEN 1 AND 5
    ),
    window_condition INTEGER NOT NULL CHECK (
        window_condition BETWEEN 1 AND 5
    ),
    noise_isolation INTEGER NOT NULL CHECK (
        noise_isolation BETWEEN 1 AND 5
    ),
    common_areas_equipment INTEGER NOT NULL CHECK (
        common_areas_equipment BETWEEN 1 AND 5
    ),
    transport_accessibility INTEGER NOT NULL CHECK (
        transport_accessibility BETWEEN 1 AND 5
    ),
    administration_quality INTEGER NOT NULL CHECK (
        administration_quality BETWEEN 1 AND 5
    ),
    residents_culture_level INTEGER NOT NULL CHECK (
        residents_culture_level BETWEEN 1 AND 5
    ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_grades_dormitory_id ON grades (dormitory_id);

CREATE INDEX IF NOT EXISTS idx_grades_user_id ON grades (user_id);

CREATE INDEX IF NOT EXISTS idx_grades_created_at ON grades (created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grades_unique_per_month ON grades (
    dormitory_id,
    user_id,
    DATE_TRUNC('month', created_at)
);