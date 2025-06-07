package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// GetAllAppointments: 管理員查詢所有預約，可依 date 篩選
func GetAllAppointments(c *gin.Context) {
	var apps []models.Appointment
	query := db.DB
	if date := c.Query("date"); date != "" {
		query = query.Where("DATE(appointment_time) = ?", date)
	}
	if err := query.Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "查詢預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "appointments": apps})
}

// UpdateAppointmentStatus: 管理員更新預約狀態 (booked, cancelled, no_show)
func UpdateAdminAppointmentStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "無效預約ID", "code": 400})
		return
	}
	var req struct { Status string `json:"status"` }
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "請提供有效狀態", "code": 400})
		return
	}
	if err := db.DB.Model(&models.Appointment{}).
		Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "更新狀態失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "狀態更新成功"})
}

// DeleteAppointment: 管理員刪除預約
func DeleteAdminAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "無效預約ID", "code": 400})
		return
	}
	if err := db.DB.Delete(&models.Appointment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "刪除預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "預約已刪除"})
}

// AssignSubstitute: 管理員指派替代醫師並通知病患
func AssignSubstitute(c *gin.Context) {
	var req struct {
		AbsentDoctorId     int    `json:"absentDoctorId"`
		SubstituteDoctorId int    `json:"substituteDoctorId"`
		Date               string `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "缺少替代資訊", "code": 400})
		return
	}
	// 更新邏輯：將所有符合條件的預約換醫師並設 status
	err := db.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND date = ? AND status = ?", req.AbsentDoctorId, req.Date, "booked").
		Updates(map[string]interface{}{"doctor_id": req.SubstituteDoctorId}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "替代指派失敗", "code": 500})
		return
	}
	// TODO: 通知病人 (簡訊/Email)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "已指派替代醫師並通知病患"})
}

// GenerateReport: 管理員報表產生器
func GenerateReport(c *gin.Context) {
	typeParam := c.Query("type")
	month := c.Query("month")
	// TODO: 根據 type、month 生產報表資料
	c.JSON(http.StatusOK, gin.H{"error": false, "reportType": typeParam, "month": month, "data": "報表示範"})
}
