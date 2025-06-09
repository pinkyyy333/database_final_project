package controllers

import (
	"clinic-backend/db"
	"clinic-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLiveQueue 取得醫師即時看診佇列與等待時間
func GetLiveQueue(c *gin.Context) {
	doctorID := c.Param("doctor_id")
	now := time.Now().UTC()
	var upcoming []models.Appointment
	db.DB.Where("doctor_id = ? AND status = ? AND appointment_time <= ?", doctorID, "checked_in", now).
		Order("appointment_time asc").
		Find(&upcoming)
	// 計算等待人數與預估等待時間
	count := len(upcoming)
	var eta time.Duration
	if count > 0 {
		// 假設每人診療 15 分鐘
		eta = time.Duration((count-1)*15) * time.Minute
	}
	c.JSON(http.StatusOK, gin.H{
		"error":          false,
		"waiting_count":  count,
		"estimated_wait": eta.String(),
		"queue":          upcoming,
	})
}
