// controllers/feedback_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// POST /api/v1/feedbacks
func CreateFeedback(c *gin.Context) {
	var req struct {
		AppointmentID  uint32   `json:"appointment_id"`
		FeedbackRating int    `json:"feedback_rating"`
		PatientComment string `json:"patient_comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   true,
			"message": "參數錯誤",
			"code":    400,
		})
		return
	}

	f := models.Feedback{
		AppointmentID:  req.AppointmentID,
		FeedbackRating: req.FeedbackRating,
		PatientComment: req.PatientComment,
	}
	if err := db.DB.Create(&f).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "建立評價失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"error":    false,
		"feedback": f,
	})
}

// GET /api/v1/doctors/:doctor_id/feedbacks
func GetDoctorFeedbacks(c *gin.Context) {
	idStr := c.Param("id")
	did, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   true,
			"message": "doctor_id 參數錯誤，必須為數字",
			"code":    400,
		})
		return
	}

	var feedbacks []models.Feedback
	if err := db.DB.
		Joins("JOIN appointments ON appointments.id = feedbacks.appointment_id").
		Where("appointments.doctor_id = ?", did).
		Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "取得評價失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"error":     false,
		"feedbacks": feedbacks,
	})
}
