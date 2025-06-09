package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "clinic-backend/db"
    "clinic-backend/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/datatypes"
)

// PatientRequest 綁定註冊請求，包含前端所有欄位
type PatientRequest struct {
	PatientID         string   `json:"patient_id"`
	PatientName       string   `json:"patient_name"`
	PatientGender     string   `json:"patient_gender"`
	PatientBirth      string   `json:"patient_birth"`
	PatientPhone      string   `json:"patient_phone"`
	Password          string   `json:"password"`
	Address           string   `json:"address"`
	EmergencyName     string   `json:"emergency_name"`
	EmergencyPhone    string   `json:"emergency_phone"`
	EmergencyRelation string   `json:"emergency_relation"`
	DrugAllergy       []string `json:"drug_allergy"`
	FoodAllergy       []string `json:"food_allergy"`
	MedicalHistory    []string `json:"medical_history"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤: " + err.Error(), "code": 400})
		return
	}

	// 檢查是否已存在
	var existing models.Patient
	if err := db.DB.First(&existing, "patient_id = ?", req.PatientID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": true, "message": "此身分證字號已被註冊", "code": 409})
		return
	}

	// 密碼雜湊
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "無法加密密碼", "code": 500})
		return
	}

	// 將 slice 轉為 JSON 格式存入 datatypes.JSON
	drugJSON, _ := json.Marshal(req.DrugAllergy)
	foodJSON, _ := json.Marshal(req.FoodAllergy)
	historyJSON, _ := json.Marshal(req.MedicalHistory)

	p := models.Patient{
		PatientID:         req.PatientID,
		PatientName:       req.PatientName,
		PatientGender:     req.PatientGender,
		PatientBirth:      req.PatientBirth,
		PatientPhone:      req.PatientPhone,
		Password:          string(hashed),
		Address:           req.Address,
		EmergencyName:     req.EmergencyName,
		EmergencyPhone:    req.EmergencyPhone,
		EmergencyRelation: req.EmergencyRelation,
		DrugAllergy:       datatypes.JSON(drugJSON),
		FoodAllergy:       datatypes.JSON(foodJSON),
		MedicalHistory:    datatypes.JSON(historyJSON),
	}

	if err := db.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "註冊失敗: " + err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "註冊成功"})
}

// LoginPatient 病患登入
func LoginPatient(c *gin.Context) {
    var req struct {
		PatientID string `json:"patient_id"`
        Password  string `json:"password"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤: " + err.Error(), "code": 400})
        return
    }

    fmt.Printf("[DEBUG] Login 嘗試 patient_id=%q\n", req.PatientID)

    // *只查 password 欄位，避免一次掃描整個 models.Patient 的 JSON 欄位導致失敗*
    var storedHash string
    if err := db.DB.Model(&models.Patient{}).
        Select("password").
        Where("patient_id = ?", req.PatientID).
        Scan(&storedHash).Error; err != nil {
        fmt.Printf("[DEBUG] 查詢密碼失敗: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "帳號或密碼錯誤", "code": 401})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
        fmt.Printf("[DEBUG] 密碼比對失敗: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "帳號或密碼錯誤", "code": 401})
        return
    }

    // 簽發 JWT，下面流程不變…
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "伺服器未設定 JWT_SECRET", "code": 500})
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Token 產生失敗", "code": 500})
        return
    }
    c.JSON(http.StatusOK, gin.H{"error": false, "token": token})
}

// GetPatientProfile 取得個人資料
func GetPatientProfile(c *gin.Context) {
	pid, _ := c.Get("patient_id")
	var p models.Patient
	if err := db.DB.First(&p, "patient_id = ?", pid.(string)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "讀取資料失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "patient": p})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	if err := db.DB.Model(&models.Patient{}).
		Where("patient_id = ?", pid.(string)).
		Updates(models.Patient{
			PatientName:   req.PatientName,
			PatientGender: req.PatientGender,
			PatientBirth:  req.PatientBirth,
			PatientPhone:  req.PatientPhone,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "更新失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "更新成功"})
}

// ChangePatientPassword 更新密碼
func ChangePatientPassword(c *gin.Context) {
	pid, _ := c.Get("patient_id")
	var req struct{ OldPassword, NewPassword string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "參數錯誤", "code": 400})
		return
	}
	var p models.Patient
	if err := db.DB.First(&p, "patient_id = ?", pid.(string)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "讀取舊密碼失敗", "code": 500})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(req.OldPassword)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "舊密碼不正確", "code": 401})
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := db.DB.Model(&models.Patient{}).
		Where("patient_id = ?", pid.(string)).
		Update("password", string(newHash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "更新密碼失敗", "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "密碼更新成功"})
}
