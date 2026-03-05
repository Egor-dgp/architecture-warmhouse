package repository

import (
    "database/sql"
    "smart-temperature-system/internal/models"
)

type SensorRepository interface {
    GetAll() ([]models.Sensor, error)
    GetByID(id int) (*models.Sensor, error)
    Create(sensor *models.CreateSensorRequest) (*models.Sensor, error)
    UpdateStatus(id int, isActive bool) error
    Delete(id int) error
    UpdateLastReading(id int, temperature float64) error
}

type sensorRepository struct {
    db *sql.DB
}

func NewSensorRepository(db *sql.DB) SensorRepository {
    return &sensorRepository{db: db}
}

func (r *sensorRepository) GetAll() ([]models.Sensor, error) {
    // SQL запрос для получения всех датчиков
    query := `
        SELECT id, name, location, is_active, min_temp, max_temp, 
               last_reading, created_at, updated_at 
        FROM sensors 
        ORDER BY id`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var sensors []models.Sensor
    for rows.Next() {
        var sensor models.Sensor
        err := rows.Scan(
            &sensor.ID,
            &sensor.Name,
            &sensor.Location,
            &sensor.IsActive,
            &sensor.MinTemp,
            &sensor.MaxTemp,
            &sensor.LastReading,
            &sensor.CreatedAt,
            &sensor.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        sensors = append(sensors, sensor)
    }
    
    return sensors, nil
}

func (r *sensorRepository) GetByID(id int) (*models.Sensor, error) {
    // SQL запрос для получения датчика по ID
    query := `
        SELECT id, name, location, is_active, min_temp, max_temp, 
               last_reading, created_at, updated_at 
        FROM sensors 
        WHERE id = $1`
    
    var sensor models.Sensor
    err := r.db.QueryRow(query, id).Scan(
        &sensor.ID,
        &sensor.Name,
        &sensor.Location,
        &sensor.IsActive,
        &sensor.MinTemp,
        &sensor.MaxTemp,
        &sensor.LastReading,
        &sensor.CreatedAt,
        &sensor.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Датчик не найден
        }
        return nil, err
    }
    
    return &sensor, nil
}

func (r *sensorRepository) Create(sensorReq *models.CreateSensorRequest) (*models.Sensor, error) {
    // SQL запрос для создания нового датчика
    query := `
        INSERT INTO sensors (name, location, min_temp, max_temp, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, true, NOW(), NOW())
        RETURNING id, name, location, is_active, min_temp, max_temp, last_reading, created_at, updated_at`
    
    var sensor models.Sensor
    err := r.db.QueryRow(
        query,
        sensorReq.Name,
        sensorReq.Location,
        sensorReq.MinTemp,
        sensorReq.MaxTemp,
    ).Scan(
        &sensor.ID,
        &sensor.Name,
        &sensor.Location,
        &sensor.IsActive,
        &sensor.MinTemp,
        &sensor.MaxTemp,
        &sensor.LastReading,
        &sensor.CreatedAt,
        &sensor.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &sensor, nil
}

func (r *sensorRepository) UpdateStatus(id int, isActive bool) error {
    // SQL запрос для обновления статуса датчика
    query := `UPDATE sensors SET is_active = $1, updated_at = NOW() WHERE id = $2`
    
    _, err := r.db.Exec(query, isActive, id)
    return err
}

func (r *sensorRepository) Delete(id int) error {
    // SQL запрос для удаления датчика
    query := `DELETE FROM sensors WHERE id = $1`
    
    _, err := r.db.Exec(query, id)
    return err
}

func (r *sensorRepository) UpdateLastReading(id int, temperature float64) error {
    // SQL запрос для обновления последнего показания температуры
    query := `UPDATE sensors SET last_reading = $1, updated_at = NOW() WHERE id = $2`
    
    _, err := r.db.Exec(query, temperature, id)
    return err
}