-- Создаем таблицу датчиков
CREATE TABLE IF NOT EXISTS sensors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(200) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    min_temp DECIMAL(5,2) DEFAULT -20.00,
    max_temp DECIMAL(5,2) DEFAULT 40.00,
    last_reading DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для быстрого поиска по имени
CREATE INDEX idx_sensors_name ON sensors(name);

-- Создаем индекс для быстрого поиска по местоположению
CREATE INDEX idx_sensors_location ON sensors(location);

-- Вставляем тестовые данные
INSERT INTO sensors (name, location, min_temp, max_temp) VALUES
    ('Kitchen Sensor', 'Kitchen', 15.0, 25.0),
    ('Living Room Sensor', 'Living Room', 18.0, 24.0),
    ('Bedroom Sensor', 'Bedroom', 16.0, 22.0),
    ('Garage Sensor', 'Garage', -10.0, 30.0),
    ('Basement Sensor', 'Basement', 10.0, 18.0)
ON CONFLICT DO NOTHING;

-- Создаем функцию для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создаем триггер для автоматического обновления updated_at
CREATE TRIGGER update_sensors_updated_at
    BEFORE UPDATE ON sensors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();