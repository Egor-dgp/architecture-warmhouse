package db // пакет db отвечает за подключение и операции с БД

import (
	"context" // для передачи контекста в запросы
	"fmt"     // для форматированных ошибок
	"time"    // для работы со временем

	"smarthome/apps/smart_home/models" // локальный импорт моделей

	"github.com/jackc/pgx/v5/pgxpool" // pgxpool — пул соединений для Postgres
)

// DB — обёртка вокруг пула соединений
type DB struct {
	Pool *pgxpool.Pool // пул соединений к Postgres
}

// New — создаёт новое подключение к БД по строке подключения connString
func New(connString string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // контекст с таймаутом на подключение
	defer cancel()                                                          // отложенное освобождение контекста

	pool, err := pgxpool.New(ctx, connString) // создаём пул соединений
	if err != nil {
		return nil, fmt.Errorf("unable to create pool: %w", err) // возвращаем ошибку, если не получилось
	}

	// проверяем соединение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()                                         // закрываем пул при ошибке
		return nil, fmt.Errorf("unable to ping db: %w", err) // возвращаем ошибку
	}

	return &DB{Pool: pool}, nil // возвращаем структуру DB с пулом
}

// Close — закрывает пул соединений
func (db *DB) Close() {
	if db.Pool != nil { // если пул не nil
		db.Pool.Close() // закрываем пул
	}
}

// CreateSensor — создаёт новый сенсор и возвращает его
func (db *DB) CreateSensor(ctx context.Context, name string) (models.Sensor, error) {
	query := `INSERT INTO sensors (name, created_at) VALUES ($1, $2) RETURNING id, name, created_at` // SQL запрос вставки
	var s models.Sensor                                                                                 // переменная для результата
	err := db.Pool.QueryRow(ctx, query, name, time.Now().UTC()).Scan(&s.ID, &s.Name, &s.CreatedAt)     // выполняем запрос и сканируем результат
	if err != nil {
		return models.Sensor{}, fmt.Errorf("create sensor: %w", err) // возвращаем ошибку, если не удалось
	}
	return s, nil // возвращаем созданный сенсор
}

// ListSensors — возвращает все сенсоры
func (db *DB) ListSensors(ctx context.Context) ([]models.Sensor, error) {
	query := `SELECT id, name, created_at FROM sensors ORDER BY id` // запрос списка сенсоров
	rows, err := db.Pool.Query(ctx, query)                          // выполняем запрос
	if err != nil {
		return nil, fmt.Errorf("list sensors: %w", err) // ошибка выполнения запроса
	}
	defer rows.Close() // не забываем закрыть rows

	var out []models.Sensor // слайс для результата
	for rows.Next() {       // итерация по строкам
		var s models.Sensor
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt); err != nil { // считываем строку
			return nil, fmt.Errorf("scan sensor: %w", err) // ошибка сканирования
		}
		out = append(out, s) // добавляем в результат
	}
	return out, nil // возвращаем список
}

// SaveMeasurement — сохраняет измерение в БД
func (db *DB) SaveMeasurement(ctx context.Context, m models.Measurement) (models.Measurement, error) {
	query := `INSERT INTO measurements (sensor_id, temperature, measured_at) VALUES ($1, $2, $3) RETURNING id, sensor_id, temperature, measured_at` // SQL вставки
	var res models.Measurement                                                                                                                                      // результат
	err := db.Pool.QueryRow(ctx, query, m.SensorID, m.Temperature, m.MeasuredAt).Scan(&res.ID, &res.SensorID, &res.Temperature, &res.MeasuredAt)                 // выполняем
	if err != nil {
		return models.Measurement{}, fmt.Errorf("save measurement: %w", err) // ошибка при сохранении
	}
	return res, nil // возвращаем сохранённое измерение
}

// GetLatestMeasurement — получает последнее измерение для сенсора
func (db *DB) GetLatestMeasurement(ctx context.Context, sensorID int) (models.Measurement, error) {
	query := `SELECT id, sensor_id, temperature, measured_at FROM measurements WHERE sensor_id = $1 ORDER BY measured_at DESC LIMIT 1` // запрос последнего измерения
	var m models.Measurement                                                                                                             // результат
	err := db.Pool.QueryRow(ctx, query, sensorID).Scan(&m.ID, &m.SensorID, &m.Temperature, &m.MeasuredAt)                                 // выполняем запрос
	if err != nil {
		return models.Measurement{}, fmt.Errorf("get latest: %w", err) // ошибка, если не найдено или другая
	}
	return m, nil // возвращаем найденное измерение
}