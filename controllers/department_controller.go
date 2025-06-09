// controllers/department_controller.go
package controllers

import (
	"net/http"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/departments
func GetAllDepartments(c *gin.Context) {
	var deps []models.Department
	if err := db.DB.Find(&deps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取得科別失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "departments": deps})
}

// GET /api/v1/departments/:dept_id/doctors
func GetDoctorsByDepartment(c *gin.Context) {
	deptID := c.Param("dept_id")
	var docs []models.Doctor
	if err := db.DB.Where("dept_id = ?", deptID).Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取得醫師失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "doctors": docs})
}
