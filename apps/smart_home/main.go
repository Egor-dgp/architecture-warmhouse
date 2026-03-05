package main // основной пакет исполняемого файла

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"smart_home/db"      // измените с "smart_home/apps/smart_home/db"
	"smart_home/handlers" // измените с "smart_home/apps/smart_home/handlers"
	"smart_home/services" // измените с "smart_home/apps/smart_home/services"
)

func main() {
	// Получаем строку подключения из окружения или используем дефолт
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// локальный дефолт для docker-compose (имя сервиса smart_db)
		connStr = "postgres://smart_user:smart_pass@smart_db:5432/smart_home_db?sslmode=disable"
	}

	// Создаём подключение к БД с контекстом и таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbConn, err := db.NewWithContext(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err) // фатал при ошибке подключения
	}
	defer dbConn.Close() // гарантированно закрываем при завершении

	// Инициализируем сервис и обработчики
	tempSvc := services.NewTemperatureService(dbConn)       // сервис генерации температур
	handlers := handlers.NewHandlers(tempSvc, dbConn)       // HTTP-обработчики

	// Настраиваем маршруты
	r := mux.NewRouter()                             // создаём роутер
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { // простой healthcheck
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")

	r.HandleFunc("/sensors", handlers.CreateSensorHandler).Methods("POST")          // создать сенсор
	r.HandleFunc("/sensors", handlers.ListSensorsHandler).Methods("GET")           // список сенсоров
	r.HandleFunc("/sensors/{id}/measure", handlers.MeasureHandler).Methods("POST") // сгенерировать и сохранить измерение
	r.HandleFunc("/sensors/{id}/latest", handlers.GetLatestHandler).Methods("GET") // получить последнее измерение

	// Порт сервера (можно переопределить через PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Запускаем HTTP сервер
	srv := &http.Server{
		Handler:      r,                // обработчик
		Addr:         ":" + port,       // адрес прослушивания
		ReadTimeout:  5 * time.Second,  // таймаут чтения
		WriteTimeout: 10 * time.Second, // таймаут записи
	}

	log.Printf("listening on %s", srv.Addr) // лог запуска
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err) // лог ошибки сервера
	}
}