package service

import (
    "math/rand"
    "smart-temperature-system/internal/models"
    "smart-temperature-system/internal/repository"
    "time"
)

type SensorService interface {
    GetAllSensors() ([]models.Sensor, error)
    GetSensor(id int) (*models.Sensor, error)
    CreateSensor(sensorReq *models.CreateSensorRequest) (*models.Sensor, error)
    ToggleSensor(id int) (*models.Sensor, error)
    DeleteSensor(id int) error
    GetTemperature(id int) (*models.TemperatureResponse, error)
}

type sensorService struct {
    repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
    return &sensorService{repo: repo}
}

func (s *sensorService) GetAllSensors() ([]models.Sensor, error) {
    return s.repo.GetAll()
}

func (s *sensorService) GetSensor(id int) (*models.Sensor, error) {
    return s.repo.GetByID(id)
}

func (s *sensorService) CreateSensor(sensorReq *models.CreateSensorRequest) (*models.Sensor, error) {
    // Устанавливаем значения по умолчанию для температур
    if sensorReq.MinTemp == 0 {
        sensorReq.MinTemp = -20.0
    }
    if sensorReq.MaxTemp == 0 {
        sensorReq.MaxTemp = 40.0
    }
    
    return s.repo.Create(sensorReq)
}

func (s *sensorService) ToggleSensor(id int) (*models.Sensor, error) {
    // Получаем текущий датчик
    sensor, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if sensor == nil {
        return nil, nil // Датчик не найден
    }
    
    // Меняем статус на противоположный
    newStatus := !sensor.IsActive
    err = s.repo.UpdateStatus(id, newStatus)
    if err != nil {
        return nil, err
    }
    
    // Получаем обновленный датчик
    return s.repo.GetByID(id)
}

func (s *sensorService) DeleteSensor(id int) error {
    return s.repo.Delete(id)
}

func (s *sensorService) GetTemperature(id int) (*models.TemperatureResponse, error) {
    // Получаем информацию о датчике
    sensor, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if sensor == nil {
        return nil, nil // Датчик не найден
    }
    
    // Если датчик выключен, возвращаем последнее показание или 0
    if !sensor.IsActive {
        return &models.TemperatureResponse{
            SensorID:    sensor.ID,
            Temperature: sensor.LastReading,
            Timestamp:   time.Now(),
            IsActive:    false,
        }, nil
    }
    
    // Генерируем случайную температуру в пределах диапазона датчика
    rand.Seed(time.Now().UnixNano())
    temperature := sensor.MinTemp + rand.Float64()*(sensor.MaxTemp-sensor.MinTemp)
    
    // Округляем до одного знака после запятой
    temperature = float64(int(temperature*10)) / 10
    
    // Обновляем последнее показание в БД
    err = s.repo.UpdateLastReading(id, temperature)
    if err != nil {
        return nil, err
    }
    
    return &models.TemperatureResponse{
        SensorID:    sensor.ID,
        Temperature: temperature,
        Timestamp:   time.Now(),
        IsActive:    true,
    }, nil
}