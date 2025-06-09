package controllers

import (
	"net/http"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

type AppointmentRequest struct {
	DepartmentID    uint      `json:"department_id"`
	DoctorID        uint      `json:"doctor_id"`
	PatientID       string    `json:"patient_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	ServiceType     string    `json:"service_type"`
}

func CreateAppointment(c *gin.Context) {
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if req.AppointmentTime.Before(time.Now().UTC()) {
		RespondError(c, http.StatusBadRequest, "預約時間必須在未來")
		return
	}
	// 衝突檢查
	var cnt int64
	if err := db.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND appointment_time = ?", req.DoctorID, req.AppointmentTime).
		Count(&cnt).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "查詢衝突失敗")
		return
	}
	if cnt > 0 {
		RespondError(c, http.StatusConflict, "時段已被預約")
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
		RespondError(c, http.StatusInternalServerError, "建立預約失敗")
		return
	}
	RespondCreated(c, gin.H{"appointment": a})
}

func GetPatientAppointments(c *gin.Context) {
	pid := c.Param("patient_id")
	var list []models.Appointment
	if err := db.DB.Where("patient_id = ?", pid).Find(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取得預約失敗")
		return
	}
	RespondOK(c, gin.H{"appointments": list})
}

func GetDoctorAppointments(c *gin.Context) {
	did := c.Param("doctor_id")
	var list []models.Appointment
	if err := db.DB.Where("doctor_id = ?", did).Find(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取得預約失敗")
		return
	}
	RespondOK(c, gin.H{"appointments": list})
}

func GetAvailableSlots(c *gin.Context) {
	doctorID := c.Param("doctor_id")
	dateStr := c.Query("date")
	if dateStr == "" {
		RespondError(c, http.StatusBadRequest, "請提供 date 參數，格式 YYYY-MM-DD")
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "date 格式錯誤，需 YYYY-MM-DD")
		return
	}
	var booked []models.Appointment
	if err := db.DB.
		Where("doctor_id = ? AND DATE(appointment_time) = ?", doctorID, dateStr).
		Find(&booked).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "查詢已預約時段失敗")
		return
	}
	occupied := make(map[string]bool)
	for _, appt := range booked {
		occupied[appt.AppointmentTime.UTC().Format("15:04")] = true
	}
	var slots []string
	for hour := 9; hour < 17; hour++ {
		t := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)
		key := t.Format("15:04")
		if !occupied[key] {
			slots = append(slots, t.Format(time.RFC3339))
		}
	}
	RespondOK(c, gin.H{"available_slots": slots})
}

func UpdateAppointmentStatus(c *gin.Context) {
	id := c.Param("appointment_id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "參數錯誤")
		return
	}
	var a models.Appointment
	if err := db.DB.First(&a, id).Error; err != nil {
		RespondError(c, http.StatusNotFound, "找不到該預約")
		return
	}
	a.Status = req.Status
	if err := db.DB.Save(&a).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "更新狀態失敗")
		return
	}
	RespondOK(c, gin.H{"appointment": a})
}

func CancelAppointment(c *gin.Context) {
	id := c.Param("appointment_id")
	var a models.Appointment
	if err := db.DB.First(&a, id).Error; err != nil {
		RespondError(c, http.StatusNotFound, "找不到該預約")
		return
	}
	a.Status = "cancelled"
	if err := db.DB.Save(&a).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取消預約失敗")
		return
	}
	RespondOK(c, gin.H{"message": "預約已取消"})
}

func CheckInAppointment(c *gin.Context) {
	id := c.Param("appointment_id")
	var appt models.Appointment
	if err := db.DB.First(&appt, id).Error; err != nil {
		RespondError(c, http.StatusNotFound, "找不到該預約")
		return
	}
	now := time.Now().UTC()
	appt.CheckInTime = &now
	appt.Status = "checked_in"
	if err := db.DB.Save(&appt).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "報到失敗")
		return
	}
	RespondOK(c, gin.H{"appointment": appt})
}
