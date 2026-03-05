package models // пакет с моделями — данные приложения

import "time" // импортируем время для меток времени

// Sensor — модель сенсора в БД/в приложении
type Sensor struct {
	ID        int       `json:"id"`         // уникальный идентификатор сенсора
	Name      string    `json:"name"`       // имя сенсора (например "kitchen")
	CreatedAt time.Time `json:"created_at"` // время создания записи
}

// Measurement — модель показания сенсора
type Measurement struct {
	ID         int       `json:"id"`          // уникальный идентификатор измерения
	SensorID   int       `json:"sensor_id"`   // внешний ключ на сенсор
	Temperature float64  `json:"temperature"` // значение температуры
	MeasuredAt time.Time `json:"measured_at"` // время измерения
}