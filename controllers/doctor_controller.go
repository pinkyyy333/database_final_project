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

// DoctorClaims 醫師用的 token payload
// 使用 jwt.RegisteredClaims 支援 v4 欄位
type DoctorClaims struct {
	DoctorID uint   `json:"doctor_id"`
	Role     string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// CreateDoctor 註冊新醫師（含密碼雜湊）
func CreateDoctor(c *gin.Context) {
	var req struct {
		DoctorID   uint32    `json:"doctor_id,string"`
		DeptID     uint32    `json:"dept_id"`
		DoctorName string    `json:"doctor_name"`
		DoctorInfo string    `json:"doctor_info"`
		Password   string    `json:"password"`
		Gender     string    `gorm:"column:gender" json:"gender"`
		Edu        string    `gorm:"column:edu"    json:"edu"`
		HireDate   time.Time `gorm:"column:hire_date" json:"hire_date"`
		Phone      string    `gorm:"column:phone"  json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "code": 400})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "雜湊密碼失敗", "code": 500})
		return
	}
	doc := models.Doctor{
		DoctorID:   req.DoctorID,
		DeptID:     req.DeptID,
		DoctorName: req.DoctorName,
		DoctorInfo: req.DoctorInfo,
		Password:   string(hashed),
		Gender:     req.Gender,
		Edu:        req.Edu,
		HireDate:   req.HireDate,
		Phone:      req.Phone,
	}
	if err := db.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "doctor": doc})
}

// GetAllDoctors 取得醫師列表（供前端預約使用，可依 dept 查詢）
func GetAllDoctors(c *gin.Context) {
	dept := c.Query("dept")
	var docs []models.Doctor
	var err error
	if dept != "" {
		err = db.DB.Where("dept_id = ?", dept).Find(&docs).Error
	} else {
		err = db.DB.Find(&docs).Error
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "取得醫生列表失敗", "code": 500})
		return
	}
	resp := make([]map[string]interface{}, len(docs))
	for i, d := range docs {
		resp[i] = map[string]interface{}{
			"id":   d.DoctorID,
			"name": d.DoctorName,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateDoctor 更新醫師資料（管理端）
func UpdateDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的ID", "code": 400})
		return
	}
	var doc models.Doctor
	if err := db.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "找不到醫師", "code": 404})
		return
	}
	var input struct {
		DeptID     uint32    `json:"dept_id"`
		DoctorName string    `json:"doctor_name"`
		DoctorInfo string    `json:"doctor_info"`
		Gender     string    `gorm:"column:gender" json:"gender"`
		Edu        string    `gorm:"column:edu"    json:"edu"`
		HireDate   time.Time `gorm:"column:hire_date" json:"hire_date"`
		Phone      string    `gorm:"column:phone"  json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "code": 400})
		return
	}
	doc.DeptID = input.DeptID
	doc.DoctorName = input.DoctorName
	doc.DoctorInfo = input.DoctorInfo
	doc.Gender = input.Gender
	doc.Edu = input.Edu
	doc.HireDate = input.HireDate
	doc.Phone = input.Phone
	if err := db.DB.Save(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "doctor": doc})
}

// DeleteDoctor 刪除醫師（管理端）
func DeleteDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "無效的ID", "code": 400})
		return
	}
	if err := db.DB.Delete(&models.Doctor{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "刪除成功"})
}

// GetDoctorSchedule 取得某醫師所有排程
func GetDoctorSchedule(c *gin.Context) {
	id := c.Param("id")
	var list []models.Appointment
	if err := db.DB.Where("doctor_id = ?", id).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "取得排程失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "schedule": list})
}

// GetPatientRecords 取得某醫師對某病患的所有病歷
func GetPatientRecords(c *gin.Context) {
	docID := c.Param("id")
	patID := c.Param("patient_id")
	var recs []models.Appointment
	if err := db.DB.
		Where("doctor_id = ? AND patient_id = ?", docID, patID).
		Find(&recs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "取得病歷失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "records": recs})
}

// LoginDoctor 醫師登入並回傳 JWT
func LoginDoctor(c *gin.Context) {
	var req struct {
		DoctorID int    `json:"doctor_id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "參數錯誤", "code": 400})
		return
	}
	var doc models.Doctor
	if err := db.DB.First(&doc, "doctor_id = ?", req.DoctorID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "帳號或密碼錯誤", "code": 401})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(doc.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "帳號或密碼錯誤", "code": 401})
		return
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "伺服器未設定 JWT_SECRET", "code": 500})
		return
	}
	claims := DoctorClaims{
		DoctorID: uint(req.DoctorID),
		Role:     "doctor",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "clinic-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Token 產生失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "token": signed})
}
