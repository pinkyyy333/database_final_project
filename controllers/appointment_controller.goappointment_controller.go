package controllers

import (
	"net/http"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// AppointmentRequest 包含服務類型
type AppointmentRequest struct {
	DepartmentID    uint      `json:"department_id"`
	DoctorID        uint      `json:"doctor_id"`
	PatientID       string    `json:"patient_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	ServiceType     string    `json:"service_type"`
}

// CreateAppointment 新增預約（含非看診服務）
func CreateAppointment(c *gin.Context) {
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid JSON", "code": 400})
		return
	}
	if req.AppointmentTime.Before(time.Now().UTC()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "預約時間必須在未來", "code": 400})
		return
	}
	// 衝突檢查
	var cnt int64
	db.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND appointment_time = ?", req.DoctorID, req.AppointmentTime).
		Count(&cnt)
	if cnt > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": true, "message": "時段已被預約", "code": 409})
		return
	}
	a := models.Appointment{
		DepartmentID:    req.DepartmentID,
		DoctorID:        req.DoctorID,
		PatientID:       req.PatientID,
		AppointmentTime: req.AppointmentTime,
		Status:          "booked",
		ServiceType:     req.ServiceType,
	}
	if err := db.DB.Create(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "建立預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"error": false, "appointment": a})
}

// GetPatientAppointments 查詢病患所有預約
func GetPatientAppointments(c *gin.Context) {
	pid := c.Param("patient_id")
	var list []models.Appointment
	if err := db.DB.Where("patient_id = ?", pid).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取得預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "appointments": list})
}

// GetDoctorAppointments 查詢醫師所有預約
func GetDoctorAppointments(c *gin.Context) {
	did := c.Param("doctor_id")
	var list []models.Appointment
	if err := db.DB.Where("doctor_id = ?", did).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取得預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "appointments": list})
}

// GetAvailableSlots 查詢醫師在指定日期的可預約時段
func GetAvailableSlots(c *gin.Context) {
	doctorID := c.Param("doctor_id")
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "請提供 date 參數，格式 YYYY-MM-DD", "code": 400})
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "date 格式錯誤，需 YYYY-MM-DD", "code": 400})
		return
	}
	var booked []models.Appointment
	if err := db.DB.Where("doctor_id = ? AND DATE(appointment_time) = ?", doctorID, dateStr).Find(&booked).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "查詢已預約時段失敗", "code": 500})
		return
	}
	occupied := make(map[string]bool)
	for _, appt := range booked {
		key := appt.AppointmentTime.UTC().Format("15:04")
		occupied[key] = true
	}
	var slots []string
	for hour := 9; hour < 17; hour++ {
		t := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)
		key := t.Format("15:04")
		if !occupied[key] {
			slots = append(slots, t.Format(time.RFC3339))
		}
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "available_slots": slots})
}

// UpdateAppointmentStatus 更新預約狀態
func UpdateAppointmentStatus(c *gin.Context) {
	id := c.Param("appointment_id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	var a models.Appointment
	if err := db.DB.First(&a, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "找不到該預約", "code": 404})
		return
	}
	a.Status = req.Status
	if err := db.DB.Save(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "更新狀態失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "appointment": a})
}

// CancelAppointment 取消預約（標記 cancelled）
func CancelAppointment(c *gin.Context) {
	id := c.Param("appointment_id")
	var a models.Appointment
	if err := db.DB.First(&a, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "找不到該預約", "code": 404})
		return
	}
	a.Status = "cancelled"
	if err := db.DB.Save(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取消預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "預約已取消"})
}

// CheckInAppointment 病患報到
func CheckInAppointment(c *gin.Context) {
	id := c.Param("appointment_id")
	var appt models.Appointment
	if err := db.DB.First(&appt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "找不到該預約", "code": 404})
		return
	}
	now := time.Now().UTC()
	appt.CheckInTime = &now
	appt.Status = "checked_in"
	if err := db.DB.Save(&appt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "報到失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "appointment": appt})
}
