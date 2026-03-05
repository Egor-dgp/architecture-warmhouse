-- Создаем базу и таблицы для проекта

CREATE TABLE IF NOT EXISTS sensors (
  id SERIAL PRIMARY KEY,           -- уникальный идентификатор
  name TEXT NOT NULL,              -- имя сенсора
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now() -- время создания
);

CREATE TABLE IF NOT EXISTS measurements (
  id SERIAL PRIMARY KEY,               -- уникальный идентификатор измерения
  sensor_id INT NOT NULL REFERENCES sensors(id) ON DELETE CASCADE, -- внешний ключ на sensor
  temperature NUMERIC(6,2) NOT NULL,  -- значение температуры с двумя знаками после запятой
  measured_at TIMESTAMP WITH TIME ZONE DEFAULT now() -- время измерения
);