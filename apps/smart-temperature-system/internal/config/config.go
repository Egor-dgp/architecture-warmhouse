package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    ServerPort  string
}

func Load() *Config {
    // Загружаем переменные окружения из .env файла
    _ = godotenv.Load()
    
    // Получаем URL базы данных из переменных окружения
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        // Значение по умолчанию для Docker Compose
        databaseURL = "postgres://postgres:password@db:5432/temperature_db?sslmode=disable"
    }
    
    // Получаем порт сервера
    serverPort := os.Getenv("SERVER_PORT")
    if serverPort == "" {
        serverPort = "8081"
    }
    
    return &Config{
        DatabaseURL: databaseURL,
        ServerPort:  serverPort,
    }
}