package handlers // пакет с HTTP-обработчиками

import (
	"context"       // для контекста
	"encoding/json" // кодирование/декодирование JSON
	"net/http"      // HTTP
	"strconv"       // конвертация строк в числа
	"time"          // время для контекстов

	"github.com/gorilla/mux" // роутер для параметров пути

	"smarthome/apps/smart_home/models"   // модели
	"smarthome/apps/smart_home/services" // сервисы
	"smarthome/apps/smart_home/db"       // доступ к БД (для ошибок и т.п.)
)

// Handlers — содержит зависимости для HTTP-обработчиков
type Handlers struct {
	svc *services.TemperatureService // сервис генерации температур
	db  *db.DB                       // доступ к БД (для прямых вызовов, если нужно)
}

// NewHandlers — конструктор обработчиков
func NewHandlers(svc *services.TemperatureService, dbConn *db.DB) *Handlers {
	return &Handlers{svc: svc, db: dbConn}
}

// CreateSensorHandler — POST /sensors создаёт новый сенсор
func (h *Handlers) CreateSensorHandler(w http.ResponseWriter, r *http.Request) {
	// Таймаут на операцию
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Входной JSON должен содержать {"name":"..."}
	var req struct {
		Name string `json:"name"` // имя сенсора
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // декодируем тело запроса
		http.Error(w, "invalid request body", http.StatusBadRequest) // если ошибка — 400
		return
	}

	// Создаём сенсор в БД напрямую через db
	s, err := h.db.CreateSensor(ctx, req.Name) // используем метод CreateSensor
	if err != nil {
		http.Error(w, "unable to create sensor", http.StatusInternalServerError) // 500 при ошибке
		return
	}

	// Возвращаем созданный объект и код 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(s)
}

// ListSensorsHandler — GET /sensors возвращает все сенсоры
func (h *Handlers) ListSensorsHandler(w http.ResponseWriter, r *http.Request) {
	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Получаем список сенсоров
	list, err := h.db.ListSensors(ctx)
	if err != nil {
		http.Error(w, "unable to list sensors", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

// MeasureHandler — POST /sensors/{id}/measure эмулирует измерение и сохраняет его
func (h *Handlers) MeasureHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем id из пути
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Парсим id в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid sensor id", http.StatusBadRequest)
		return
	}

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Выполняем измерение и сохранение через сервис
	m, err := h.svc.MeasureAndSave(ctx, id)
	if err != nil {
		http.Error(w, "unable to measure or save", http.StatusInternalServerError)
		return
	}

	// Возвращаем сохранённое измерение
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(m)
}

// GetLatestHandler — GET /sensors/{id}/latest возвращает последнее измерение
func (h *Handlers) GetLatestHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем id из пути
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr) // парсим
	if err != nil {
		http.Error(w, "invalid sensor id", http.StatusBadRequest)
		return
	}

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Получаем последнее измерение из БД
	m, err := h.db.GetLatestMeasurement(ctx, id)
	if err != nil {
		http.Error(w, "measurement not found", http.StatusNotFound)
		return
	}

	// Возвращаем JSON с измерением
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(m)
}