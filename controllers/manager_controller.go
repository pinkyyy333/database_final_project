package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/v1/manager/doctors  (示範 manager CRUD，實際路由按你需求)
func CreateManager(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "雜湊密碼失敗", "code": 500})
		return
	}
	m := models.Manager{Username: req.Username, Password: string(hashed)}
	if err := db.DB.Create(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "建立失敗", "code": 500})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"error": false, "manager": m})
}

func GetAllManagers(c *gin.Context) {
	var list []models.Manager
	if err := db.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "查詢失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "managers": list})
}

func UpdateManager(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "ID 錯誤", "code": 400})
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	update := map[string]interface{}{"username": req.Username}
	if req.Password != "" {
		h, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		update["password"] = string(h)
	}
	if err := db.DB.Model(&models.Manager{}).
		Where("manager_id = ?", id).
		Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "更新失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "更新成功"})
}

func DeleteManager(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "ID 錯誤", "code": 400})
		return
	}
	if err := db.DB.Delete(&models.Manager{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "刪除失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "刪除成功"})
}

// GET /api/v1/manager/reports
func GenerateReport(c *gin.Context) {
	// TODO: 實際產出報表邏輯
	c.JSON(http.StatusOK, gin.H{"error": false, "report": "報表內容示範"})
}

func GetAllAppointments(c *gin.Context) {
	var list []models.Appointment
	if err := db.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "查詢所有預約失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "appointments": list})
}

// POST /api/v1/manager/login
func LoginManager(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	var m models.Manager
	if err := db.DB.First(&m, "username = ?", req.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "帳號或密碼錯誤", "code": 401})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "帳號或密碼錯誤", "code": 401})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
