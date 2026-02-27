package services // пакет с бизнес-логикой

import (
	"context"   // контексты для операций
	"math/rand" // генератор случайных чисел
	"time"      // время для измерений и сидирования

	"smarthome/apps/smart_home/db"     // репозиторий/слой доступа к БД
	"smarthome/apps/smart_home/models" // модели
)

// TemperatureService — сервис для генерации и сохранения температур
type TemperatureService struct {
	db *db.DB // зависимость — слой доступа к БД
}

// NewTemperatureService — конструктор сервиса
func NewTemperatureService(d *db.DB) *TemperatureService {
	rand.Seed(time.Now().UnixNano()) // сеем генератор случайных чисел
	return &TemperatureService{db: d} // возвращаем новый сервис
}

// GenerateRandomTemperature — генерирует случайную температуру в заданном диапазоне с округлением
func (s *TemperatureService) GenerateRandomTemperature(min, max float64) float64 {
	// получаем случайное значение в диапазоне [min,max)
	v := min + rand.Float64()*(max-min)
	// округляем до двух знаков
	return mathRound(v, 2)
}

// MeasureAndSave — эмитирует измерение для sensorID, сохраняет в БД и возвращает измерение
func (s *TemperatureService) MeasureAndSave(ctx context.Context, sensorID int) (models.Measurement, error) {
	// Генерируем температуру, например от -10 до 40
	temp := s.GenerateRandomTemperature(-10.0, 40.0)

	// Формируем структуру измерения
	m := models.Measurement{
		SensorID:   sensorID,         // привязка к сенсору
		Temperature: temp,            // сгенерированное значение
		MeasuredAt: time.Now().UTC(), // время измерения
	}

	// Сохраняем в БД через слой db
	saved, err := s.db.SaveMeasurement(ctx, m)
	if err != nil {
		return models.Measurement{}, err // возвращаем ошибку при сохранении
	}
	return saved, nil // возвращаем сохранённое измерение
}

// mathRound — простое округление float до n знаков
func mathRound(val float64, precision int) float64 {
	pow := 1.0
	for i := 0; i < precision; i++ { // домножаем на 10^precision
		pow *= 10
	}
	// округление с добавлением 0.5 (приблизительно)
	return float64(int(val*pow+0.5)) / pow
}