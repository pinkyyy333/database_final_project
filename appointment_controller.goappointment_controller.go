package controllers

import (
	"clinic-backend/models"
	"clinic-backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentRequest struct {
	DepartmentID    int       `json:"department_id"`
	DoctorID        int       `json:"doctor_id"`
	PatientID       int       `json:"patient_id"`
	AppointmentTime time.Time `json:"appointment_time"`
}

// 預約掛號
func CreateAppointment(c *gin.Context) {
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.AppointmentTime.Before(time.Now().UTC()) {
		utils.RespondError(c, http.StatusBadRequest, "Appointment time must be in the future")
		return
	}

	// 檢查時段是否已被預約
	if models.IsTimeSlotTaken(req.DoctorID, req.AppointmentTime) {
		utils.RespondError(c, http.StatusConflict, "Time slot already booked")
		return
	}

	// 建立預約
	if err := models.CreateAppointment(req.DepartmentID, req.DoctorID, req.PatientID, req.AppointmentTime); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create appointment")
		return
	}

	utils.RespondSuccess(c, "Appointment created successfully")
}

// 病患預約清單
func GetPatientAppointments(c *gin.Context) {
	patientID := c.Param("patient_id")

	// 注意：GetAppointmentsByPatient 會回傳 ([]Appointment, error)
	appointments, err := models.GetAppointmentsByPatient(patientID)
	if err != nil {
		// 如果查詢過程中發生錯誤，回傳 500 並附上錯誤訊息
		utils.RespondError(c, http.StatusInternalServerError, "Failed to get patient appointments")
		return
	}

	// 成功就把 appointments slice 回給前端
	c.JSON(http.StatusOK, appointments)
}

// 醫師預約清單
func GetDoctorAppointments(c *gin.Context) {
	doctorID := c.Param("doctor_id")

	// 同樣要接收 error
	appointments, err := models.GetAppointmentsByDoctor(doctorID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to get doctor appointments")
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// 更新預約狀態
func UpdateAppointmentStatus(c *gin.Context) {
	appointmentID := c.Param("appointment_id")

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid status")
		return
	}

	if !models.IsValidStatus(body.Status) {
		utils.RespondError(c, http.StatusBadRequest, "Invalid status value")
		return
	}

	if err := models.UpdateAppointmentStatus(appointmentID, body.Status); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update appointment")
		return
	}

	utils.RespondSuccess(c, "Appointment status updated")
}
