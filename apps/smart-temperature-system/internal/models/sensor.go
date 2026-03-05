package models

import "time"

type Sensor struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Location    string    `json:"location"`
    IsActive    bool      `json:"is_active"`
    MinTemp     float64   `json:"min_temp"`
    MaxTemp     float64   `json:"max_temp"`
    LastReading float64   `json:"last_reading"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSensorRequest struct {
    Name     string  `json:"name" binding:"required"`
    Location string  `json:"location" binding:"required"`
    MinTemp  float64 `json:"min_temp"`
    MaxTemp  float64 `json:"max_temp"`
}

type TemperatureResponse struct {
    SensorID    int       `json:"sensor_id"`
    Temperature float64   `json:"temperature"`
    Timestamp   time.Time `json:"timestamp"`
    IsActive    bool      `json:"is_active"`
}