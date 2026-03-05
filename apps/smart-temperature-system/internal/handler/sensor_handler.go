package handler

import (
    "net/http"
    "strconv"
    "smart-temperature-system/internal/models"
    "smart-temperature-system/internal/service"
    
    "github.com/gin-gonic/gin"
)

type SensorHandler struct {
    service service.SensorService
}

func NewSensorHandler(service service.SensorService) *SensorHandler {
    return &SensorHandler{service: service}
}

func (h *SensorHandler) GetAllSensors(c *gin.Context) {
    sensors, err := h.service.GetAllSensors()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, sensors)
}

func (h *SensorHandler) GetSensor(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }
    
    sensor, err := h.service.GetSensor(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if sensor == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
        return
    }
    
    c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandler) CreateSensor(c *gin.Context) {
    var sensorReq models.CreateSensorRequest
    
    if err := c.ShouldBindJSON(&sensorReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    sensor, err := h.service.CreateSensor(&sensorReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, sensor)
}

func (h *SensorHandler) ToggleSensor(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }
    
    sensor, err := h.service.ToggleSensor(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if sensor == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Sensor status updated",
        "sensor": sensor,
    })
}

func (h *SensorHandler) DeleteSensor(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }
    
    err = h.service.DeleteSensor(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted successfully"})
}

func (h *SensorHandler) GetTemperature(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }
    
    tempResponse, err := h.service.GetTemperature(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if tempResponse == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
        return
    }
    
    c.JSON(http.StatusOK, tempResponse)
}