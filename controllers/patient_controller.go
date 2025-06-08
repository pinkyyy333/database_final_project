package controllers

import (
	"net/http"
	"os"
	"time"
	"fmt"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// 用於綁定註冊/登入請求
type PatientRequest struct {
	PatientID     string `json:"patient_id"`
	PatientName   string `json:"patient_name"`
	PatientGender string `json:"patient_gender"`
	PatientBirth  string `json:"patient_birth"`
	PatientPhone  string `json:"patient_phone"`
	Password      string `json:"password"`
}

// JWTClaims 定義 token payload
type JWTClaims struct {
	PatientID string `json:"patient_id"`
	jwt.RegisteredClaims
}

// RegisterPatient 病患註冊
func RegisterPatient(c *gin.Context) {
	var req PatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "參數錯誤: " + err.Error(), "code": 400})
		return
	}
	fmt.Printf("RegisterPatient received: %+v\n", req)

	// 檢查是否已存在
	var existing models.Patient
	if err := db.DB.First(&existing, "Patient_ID = ?", req.PatientID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "此身分證字號已被註冊", "code": 409})
		return
	}

	// 密碼雜湊
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "無法加密密碼", "code": 500})
		return
	}

	p := models.Patient{
		PatientID:     req.PatientID,
		PatientName:   req.PatientName,
		PatientGender: req.PatientGender,
		PatientBirth:  req.PatientBirth,
		PatientPhone:  req.PatientPhone,
		Password:      string(hashed),
	}
	if err := db.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "註冊失敗: " + err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "註冊成功"})
}

// LoginPatient 病患登入
func LoginPatient(c *gin.Context) {
	var req struct {
		PatientID string `json:"patient_id"`
		Password  string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "參數錯誤: " + err.Error(), "code": 400})
		return
	}

	var p models.Patient
	if err := db.DB.First(&p, "Patient_ID = ?", req.PatientID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "帳號或密碼錯誤", "code": 401})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "帳號或密碼錯誤", "code": 401})
		return
	}

	// 簽發 JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "伺服器未設定 JWT_SECRET", "code": 500})
		return
	}
	exp := time.Now().Add(24 * time.Hour)
	claims := JWTClaims{
		PatientID: req.PatientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Token 產生失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
}

// GetPatientProfile 取得個人資料
func GetPatientProfile(c *gin.Context) {
	pid, _ := c.Get("patient_id")
	var p models.Patient
	if err := db.DB.First(&p, "Patient_ID = ?", pid.(string)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "讀取資料失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "patient": p})
}

// UpdatePatientProfile 更新個人資料
func UpdatePatientProfile(c *gin.Context) {
	pid, _ := c.Get("patient_id")
	var req struct {
		PatientName   string `json:"patient_name"`
		PatientGender string `json:"patient_gender"`
		PatientBirth  string `json:"patient_birth"`
		PatientPhone  string `json:"patient_phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "參數錯誤", "code": 400})
		return
	}
	if err := db.DB.Model(&models.Patient{}).Where("Patient_ID = ?", pid.(string)).
		Updates(models.Patient{
			PatientName:   req.PatientName,
			PatientGender: req.PatientGender,
			PatientBirth:  req.PatientBirth,
			PatientPhone:  req.PatientPhone,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功"})
}

// ChangePatientPassword 更新密碼
func ChangePatientPassword(c *gin.Context) {
	pid, _ := c.Get("patient_id")
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "參數錯誤", "code": 400})
		return
	}
	var p models.Patient
	if err := db.DB.First(&p, "Patient_ID = ?", pid.(string)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "讀取舊密碼失敗", "code": 500})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.OldPassword)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "舊密碼不正確", "code": 401})
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := db.DB.Model(&models.Patient{}).Where("Patient_ID = ?", pid.(string)).Update("Password", string(newHash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新密碼失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "密碼更新成功"})
}
