package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// (原本的 CreateDoctor, GetAllDoctors, UpdateDoctor, DeleteDoctor 如前)

// GET /api/v1/doctors/:id/schedule
func GetDoctorSchedule(c *gin.Context) {
	did := c.Param("id")
	var list []models.Appointment
	if err := db.DB.Where("doctor_id = ?", did).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取得排程失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "schedule": list})
}

// GET /api/v1/doctors/:id/patients/:patient_id/records
func GetPatientRecords(c *gin.Context) {
	did := c.Param("id")
	pid := c.Param("patient_id")
	var recs []models.Appointment
	if err := db.DB.Where("doctor_id = ? AND patient_id = ?", did, pid).Find(&recs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "取得病歷失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "records": recs})
}

func CreateDoctor(c *gin.Context) {
	var doc models.Doctor
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	if err := db.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

// GET /api/v1/manager/doctors
func GetAllDoctors(c *gin.Context) {
	var docs []models.Doctor
	if err := db.DB.Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

// PUT /api/v1/manager/doctors/:id
func UpdateDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid ID"})
		return
	}
	var doc models.Doctor
	if err := db.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "Doctor not found"})
		return
	}

	var input models.Doctor
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	// 依你的 model 欄位來更新
	doc.DoctorName = input.DoctorName
	doc.DoctorInfo = input.DoctorInfo
	doc.DeptID = input.DeptID

	if err := db.DB.Save(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

// DELETE /api/v1/manager/doctors/:id
func DeleteDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid ID"})
		return
	}
	if err := db.DB.Delete(&models.Doctor{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func LoginDoctor(c *gin.Context) {
	var req struct {
		DoctorID string `json:"doctor_id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	var doc models.Doctor
	if err := db.DB.First(&doc, "doctor_id = ?", req.DoctorID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "帳號或密碼錯誤", "code": 401})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(doc.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "帳號或密碼錯誤", "code": 401})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
