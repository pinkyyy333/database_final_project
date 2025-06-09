package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/v1/manager/doctors
func CreateManager(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "參數錯誤",
			"code":    400,
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "雜湊密碼失敗",
			"code":    500,
		})
		return
	}

	m := models.Manager{
		Username: req.Username,
		Password: string(hashed),
	}
	if err := db.DB.Create(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "建立失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"manager": m,
	})
}

func GetAllManagers(c *gin.Context) {
	var list []models.Manager
	if err := db.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查詢失敗",
			"code":    500,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"managers": list,
	})
}

func UpdateManager(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID 錯誤",
			"code":    400,
		})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "參數錯誤",
			"code":    400,
		})
		return
	}

	update := map[string]interface{}{"username": req.Username}
	if req.Password != "" {
		if h, e := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); e == nil {
			update["password"] = string(h)
		}
	}
	if err := db.DB.Model(&models.Manager{}).
		Where("manager_id = ?", id).
		Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}

func DeleteManager(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID 錯誤",
			"code":    400,
		})
		return
	}
	if err := db.DB.Delete(&models.Manager{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "刪除失敗",
			"code":    500,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "刪除成功",
	})
}

// POST /api/v1/manager/login
func LoginManager(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "參數錯誤",
			"code":    400,
		})
		return
	}

	var m models.Manager
	if err := db.DB.First(&m, "username = ?", req.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "帳號或密碼錯誤",
			"code":    401,
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "帳號或密碼錯誤",
			"code":    401,
		})
		return
	}

	// 簽發 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"manager_id": m.ManagerID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Token 簽發失敗",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
	})
}

func GenerateReport(c *gin.Context) {
	var totalDoctors int64
	if err := db.DB.Model(&models.Doctor{}).Count(&totalDoctors).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "查詢醫師總數失敗")
		return
	}
	var totalAppointments int64
	if err := db.DB.Model(&models.Appointment{}).Count(&totalAppointments).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "查詢預約總數失敗")
		return
	}
	var totalFeedbacks int64
	if err := db.DB.Model(&models.Feedback{}).Count(&totalFeedbacks).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "查詢評價總數失敗")
		return
	}

	report := gin.H{
		"total_doctors":      totalDoctors,
		"total_appointments": totalAppointments,
		"total_feedbacks":    totalFeedbacks,
	}
	RespondOK(c, gin.H{"report": report})
}
