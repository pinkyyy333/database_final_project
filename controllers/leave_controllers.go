// controllers/leave_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// POST /api/v1/leaves
func CreateLeave(c *gin.Context) {
	var req models.DoctorLeave
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   true,
			"message": "參數錯誤",
			"code":    400,
		})
		return
	}
	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "建立請假紀錄失敗",
			"code":    500,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"error":   false,
		"leave":   req,
	})
}

// GET /api/v1/leaves/:leave_id/substitutes
func GetSubstitutes(c *gin.Context) {
	leaveIDStr := c.Param("leave_id")
	leaveID, err := strconv.Atoi(leaveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   true,
			"message": "leave_id 參數錯誤，必須為數字",
			"code":    400,
		})
		return
	}

	var leave models.DoctorLeave
	if err := db.DB.First(&leave, leaveID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   true,
			"message": "找不到請假紀錄",
			"code":    404,
		})
		return
	}

	var deptID int
	if err := db.DB.
		Raw("SELECT dept_id FROM doctors WHERE doctor_id = ?", leave.DoctorID).
		Scan(&deptID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "查詢科別失敗",
			"code":    500,
		})
		return
	}

	var subs []models.Doctor
	if err := db.DB.
		Where("dept_id = ? AND doctor_id != ?", deptID, leave.DoctorID).
		Find(&subs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "查詢替代醫師失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"error":       false,
		"substitutes": subs,
	})
}
