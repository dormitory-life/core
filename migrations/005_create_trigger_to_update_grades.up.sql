CREATE OR REPLACE FUNCTION update_dormitory_averages_trigger()
RETURNS TRIGGER AS $$
DECLARE
    v_period_date DATE;
BEGIN
    -- Определяем период
    v_period_date := DATE_TRUNC('month', COALESCE(NEW.created_at, CURRENT_TIMESTAMP))::DATE;
    
    -- UPSERT напрямую в триггере
    INSERT INTO dormitory_average_grades (
        dormitory_id,
        period_date,
        avg_bathroom_cleanliness,
        avg_corridor_cleanliness,
        avg_kitchen_cleanliness,
        avg_cleaning_frequency,
        avg_room_spaciousness,
        avg_corridor_spaciousness,
        avg_kitchen_spaciousness,
        avg_shower_location_convenience,
        avg_equipment_maintenance,
        avg_window_condition,
        avg_noise_isolation,
        avg_common_areas_equipment,
        avg_transport_accessibility,
        avg_administration_quality,
        avg_residents_culture_level,
        overall_average,
        total_ratings
    )
    SELECT 
        NEW.dormitory_id,
        v_period_date,
        AVG(bathroom_cleanliness::DECIMAL)::DECIMAL(3,2),
        AVG(corridor_cleanliness::DECIMAL)::DECIMAL(3,2),
        AVG(kitchen_cleanliness::DECIMAL)::DECIMAL(3,2),
        AVG(cleaning_frequency::DECIMAL)::DECIMAL(3,2),
        AVG(room_spaciousness::DECIMAL)::DECIMAL(3,2),
        AVG(corridor_spaciousness::DECIMAL)::DECIMAL(3,2),
        AVG(kitchen_spaciousness::DECIMAL)::DECIMAL(3,2),
        AVG(shower_location_convenience::DECIMAL)::DECIMAL(3,2),
        AVG(equipment_maintenance::DECIMAL)::DECIMAL(3,2),
        AVG(window_condition::DECIMAL)::DECIMAL(3,2),
        AVG(noise_isolation::DECIMAL)::DECIMAL(3,2),
        AVG(common_areas_equipment::DECIMAL)::DECIMAL(3,2),
        AVG(transport_accessibility::DECIMAL)::DECIMAL(3,2),
        AVG(administration_quality::DECIMAL)::DECIMAL(3,2),
        AVG(residents_culture_level::DECIMAL)::DECIMAL(3,2),
        AVG(
            (bathroom_cleanliness + corridor_cleanliness + kitchen_cleanliness +
             cleaning_frequency + room_spaciousness + corridor_spaciousness +
             kitchen_spaciousness + shower_location_convenience + equipment_maintenance +
             window_condition + noise_isolation + common_areas_equipment +
             transport_accessibility + administration_quality + residents_culture_level)::DECIMAL / 15
        )::DECIMAL(3,2),
        COUNT(*)
    FROM grades
    WHERE dormitory_id = NEW.dormitory_id
      AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', v_period_date)
    ON CONFLICT (dormitory_id, period_date) 
    DO UPDATE SET
        avg_bathroom_cleanliness = EXCLUDED.avg_bathroom_cleanliness,
        avg_corridor_cleanliness = EXCLUDED.avg_corridor_cleanliness,
        avg_kitchen_cleanliness = EXCLUDED.avg_kitchen_cleanliness,
        avg_cleaning_frequency = EXCLUDED.avg_cleaning_frequency,
        avg_room_spaciousness = EXCLUDED.avg_room_spaciousness,
        avg_corridor_spaciousness = EXCLUDED.avg_corridor_spaciousness,
        avg_kitchen_spaciousness = EXCLUDED.avg_kitchen_spaciousness,
        avg_shower_location_convenience = EXCLUDED.avg_shower_location_convenience,
        avg_equipment_maintenance = EXCLUDED.avg_equipment_maintenance,
        avg_window_condition = EXCLUDED.avg_window_condition,
        avg_noise_isolation = EXCLUDED.avg_noise_isolation,
        avg_common_areas_equipment = EXCLUDED.avg_common_areas_equipment,
        avg_transport_accessibility = EXCLUDED.avg_transport_accessibility,
        avg_administration_quality = EXCLUDED.avg_administration_quality,
        avg_residents_culture_level = EXCLUDED.avg_residents_culture_level,
        overall_average = EXCLUDED.overall_average,
        total_ratings = EXCLUDED.total_ratings,
        updated_at = CURRENT_TIMESTAMP;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер на INSERT (после добавления новой оценки)
CREATE TRIGGER grades_insert_trigger
AFTER INSERT ON grades
FOR EACH ROW
EXECUTE FUNCTION update_dormitory_averages_trigger();