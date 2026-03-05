package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "smart-temperature-system/internal/config"
    "smart-temperature-system/internal/database"
    "smart-temperature-system/internal/handler"
    "smart-temperature-system/internal/repository"
    "smart-temperature-system/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    // Загружаем конфигурацию
    cfg := config.Load()
    
    // Подключаемся к базе данных
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Инициализируем репозиторий
    repo := repository.NewSensorRepository(db)
    
    // Инициализируем сервис
    sensorService := service.NewSensorService(repo)
    
    // Инициализируем обработчики
    sensorHandler := handler.NewSensorHandler(sensorService)
    
    // Настраиваем Gin роутер
    router := gin.Default()
    
    // Настраиваем маршруты
    router.GET("/api/sensors", sensorHandler.GetAllSensors)
    router.GET("/api/sensors/:id/temperature", sensorHandler.GetTemperature)
    router.POST("/api/sensors/:id/toggle", sensorHandler.ToggleSensor)
    router.POST("/api/sensors", sensorHandler.CreateSensor)
    router.GET("/api/sensors/:id", sensorHandler.GetSensor)
    router.DELETE("/api/sensors/:id", sensorHandler.DeleteSensor)
    
    // Запускаем сервер в горутине
    go func() {
        if err := router.Run(":8081"); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()
    
    log.Println("Server started on port 8081")
    
    // Ожидаем сигнал завершения
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
}