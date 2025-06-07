package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAppointmentSlots returns the available appointment slots for a given doctor and month
func GetAppointmentSlots(c *gin.Context) {
	doctorID := c.Query("doctorId")
	month := c.Query("month")

	// TODO: load real data from DB instead of hardcoded sample
	c.JSON(http.StatusOK, gin.H{
		"doctorId": doctorID,
		"month":    month,
		"slots": []gin.H{
			{"date": month + "-01", "sessions": []string{"morning", "afternoon", "evening"}},
			{"date": month + "-02", "sessions": []string{"morning", "evening"}},
		},
	})
}

// UpdateAppointmentSlots updates the available slots for a doctor
func UpdateAppointmentSlots(c *gin.Context) {
	var payload struct {
		DoctorID string                `json:"doctorId"`
		Slots    []map[string][]string `json:"slots"` // each item has date and sessions
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤"})
		return
	}

	// TODO: persist payload.Slots into DB

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "可預約時段已更新"})
}
