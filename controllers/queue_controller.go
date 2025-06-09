// controllers/queue_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/doctors/:doctor_id/queue
func GetLiveQueue(c *gin.Context) {
	doctorIDStr := c.Param("doctor_id")
	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   true,
			"message": "doctor_id 參數錯誤，必須為數字",
			"code":    400,
		})
		return
	}

	now := time.Now().UTC()
	var upcoming []models.Appointment
	if err := db.DB.
		Where("doctor_id = ? AND status = ? AND appointment_time <= ?", doctorID, "checked_in", now).
		Order("appointment_time asc").
		Find(&upcoming).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "取得佇列資料失敗",
			"code":    500,
		})
		return
	}

	count := len(upcoming)
	var eta time.Duration
	if count > 0 {
		// 假設每人診療 15 分鐘
		eta = time.Duration((count-1)*15) * time.Minute
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"error":          false,
		"waiting_count":  count,
		"estimated_wait": eta.String(),
		"queue":          upcoming,
	})
}
