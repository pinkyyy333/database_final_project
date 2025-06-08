package controllers

import (
	"clinic-backend/db"
	"clinic-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateLeave 新增請假與替代
func CreateLeave(c *gin.Context) {
	var req models.DoctorLeave
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤"})
		return
	}
	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "建立請假紀錄失敗"})
		return
	}
	// 通知受影響病患
	// TODO: 實作呼叫 ReminderService 或 LINE Notify
	c.JSON(http.StatusCreated, gin.H{"error": false, "leave": req})
}

// GetSubstitutes 查詢同科可用醫師
func GetSubstitutes(c *gin.Context) {
	leaveID := c.Param("leave_id")
	var leave models.DoctorLeave
	if err := db.DB.First(&leave, leaveID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "找不到請假紀錄"})
		return
	}
	// 撈出同科所有醫師，排除請假醫師
	var deptID int
	db.DB.Raw(`SELECT dept_id FROM doctors WHERE doctor_id = ?`, leave.DoctorID).Scan(&deptID)
	var subs []models.Doctor
	db.DB.Where("dept_id = ? AND doctor_id != ?", deptID, leave.DoctorID).Find(&subs)
	c.JSON(http.StatusOK, gin.H{"error": false, "substitutes": subs})
}
