// controllers/department_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/departments
func GetAllDepartments(c *gin.Context) {
	var deps []models.Department
	if err := db.DB.Find(&deps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "取得科別失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"departments": deps,
	})
}

// GET /api/v1/departments/:dept_id/doctors
func GetDoctorsByDepartment(c *gin.Context) {
	deptIDStr := c.Param("dept_id")
	deptID, err := strconv.Atoi(deptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   true,
			"message": "dept_id 參數錯誤，必須為數字",
			"code":    400,
		})
		return
	}

	var docs []models.Doctor
	if err := db.DB.
		Where("dept_id = ?", deptID).
		Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   true,
			"message": "取得醫師失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"doctors": docs,
	})
}
