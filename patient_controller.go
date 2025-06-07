// controllers/patient_controller.go
package controllers

import (
	"clinic-backend/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Patient 註冊與登入時使用的請求結構
type Patient struct {
	PatientID     string `json:"patient_id"`     // 身分證字號
	PatientName   string `json:"patient_name"`   // 病患姓名
	PatientGender string `json:"patient_gender"` // 性別 (e.g. "M", "F")
	PatientBirth  string `json:"patient_birth"`  // 出生日期 (UTC ISO8601)
	PatientPhone  string `json:"patient_phone"`  // 聯絡電話
	Password      string `json:"password"`       // 密碼 (明文，controller 內會做 bcrypt 加密)
}

// JWTClaims - 用於產生與解析 JWT Token 的 claims 結構
// 這裡使用 github.com/golang-jwt/jwt/v4 的 RegisteredClaims
type JWTClaims struct {
	PatientID string `json:"patient_id"`
	jwt.RegisteredClaims
}

// RegisterPatient 處理病患註冊 (HTTP POST /api/v1/patients/register)
func RegisterPatient(c *gin.Context) {
	var req Patient
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "請求參數錯誤: " + err.Error(),
			"code":    400,
		})
		return
	}

	// 檢查是否已有相同 PatientID (身分證字號) 註冊
	var exists string
	err := db.DB.QueryRow("SELECT patient_id FROM patients WHERE patient_id = ?", req.PatientID).Scan(&exists)
	if err == nil {
		// 找到重複
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "此身分證字號已被註冊",
			"code":    409,
		})
		return
	}

	// 對密碼做 bcrypt 雜湊
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "伺服器錯誤，無法加密密碼",
			"code":    500,
		})
		return
	}

	// 插入資料庫
	_, err = db.DB.Exec(
		`INSERT INTO patients (
            patient_id, patient_name, patient_gender,
            patient_birth, patient_phone, password
        ) VALUES (?, ?, ?, ?, ?, ?)`,
		req.PatientID,
		req.PatientName,
		req.PatientGender,
		req.PatientBirth,
		req.PatientPhone,
		string(hashedPw),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "註冊失敗: " + err.Error(),
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "註冊成功",
	})
}

// LoginPatient 處理病患登入 (HTTP POST /api/v1/patients/login)
func LoginPatient(c *gin.Context) {
	type LoginReq struct {
		PatientID string `json:"patient_id"`
		Password  string `json:"password"`
	}
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "請求參數錯誤: " + err.Error(),
			"code":    400,
		})
		return
	}

	// 從資料庫撈出該 patient_id 的雜湊密碼
	var hashedPw string
	err := db.DB.QueryRow("SELECT password FROM patients WHERE patient_id = ?", req.PatientID).Scan(&hashedPw)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "帳號或密碼錯誤",
			"code":    401,
		})
		return
	}

	// 比對密碼
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "帳號或密碼錯誤",
			"code":    401,
		})
		return
	}

	// 產生 JWT Token，有效期限 24 小時
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "伺服器未設定 JWT_SECRET",
			"code":    500,
		})
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)

	// 建立 claims，使用 RegisteredClaims 裡的 ExpiresAt 欄位
	claims := &JWTClaims{
		PatientID: req.PatientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// 你也可以加上其他欄位，例如 Issuer、Subject 等
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "無法產生 Token",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"token": tokenString,
	})
}

// GetPatientProfile 取得目前已登入病患的個人資料 (HTTP GET /api/v1/patients/profile)
// 前提：middleware 已將 patient_id 放到 context 內，key 為 "patient_id"
func GetPatientProfile(c *gin.Context) {
	pid, exists := c.Get("patient_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "請先登入",
			"code":    401,
		})
		return
	}
	patientID := pid.(string)

	// 從資料庫撈出病患資料（不回傳 password）
	var patient struct {
		PatientID     string `json:"patient_id"`
		PatientName   string `json:"patient_name"`
		PatientGender string `json:"patient_gender"`
		PatientBirth  string `json:"patient_birth"`
		PatientPhone  string `json:"patient_phone"`
	}
	err := db.DB.QueryRow(
		`SELECT patient_id, patient_name, patient_gender, patient_birth, patient_phone
         FROM patients WHERE patient_id = ?`,
		patientID,
	).Scan(
		&patient.PatientID,
		&patient.PatientName,
		&patient.PatientGender,
		&patient.PatientBirth,
		&patient.PatientPhone,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "無法取得病患資料: " + err.Error(),
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"patient": patient,
	})
}

// UpdatePatientProfile 更新目前已登入病患的個人資料 (HTTP PUT /api/v1/patients/profile)
// 只能更新 patient_name、patient_gender、patient_birth、patient_phone；密碼需另開 endpoint 更改
func UpdatePatientProfile(c *gin.Context) {
	pid, exists := c.Get("patient_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "請先登入",
			"code":    401,
		})
		return
	}
	patientID := pid.(string)

	// 只允許更新以下欄位
	var req struct {
		PatientName   string `json:"patient_name"`
		PatientGender string `json:"patient_gender"`
		PatientBirth  string `json:"patient_birth"`
		PatientPhone  string `json:"patient_phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "請求參數錯誤: " + err.Error(),
			"code":    400,
		})
		return
	}

	// 執行更新
	_, err := db.DB.Exec(
		`UPDATE patients SET
            patient_name = ?, patient_gender = ?, patient_birth = ?, patient_phone = ?
         WHERE patient_id = ?`,
		req.PatientName,
		req.PatientGender,
		req.PatientBirth,
		req.PatientPhone,
		patientID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "更新失敗: " + err.Error(),
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "個人資料更新成功",
	})
}

// ChangePatientPassword 更改目前已登入病患的密碼 (HTTP PUT /api/v1/patients/password)
func ChangePatientPassword(c *gin.Context) {
	pid, exists := c.Get("patient_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "請先登入",
			"code":    401,
		})
		return
	}
	patientID := pid.(string)

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "請求參數錯誤: " + err.Error(),
			"code":    400,
		})
		return
	}

	// 先取出舊雜湊密碼
	var hashedPw string
	err := db.DB.QueryRow("SELECT password FROM patients WHERE patient_id = ?", patientID).Scan(&hashedPw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "無法取得舊密碼: " + err.Error(),
			"code":    500,
		})
		return
	}

	// 驗證舊密碼
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "舊密碼不正確",
			"code":    401,
		})
		return
	}

	// 新密碼做 bcrypt 雜湊
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "伺服器錯誤，無法加密新密碼",
			"code":    500,
		})
		return
	}

	// 寫回資料庫
	_, err = db.DB.Exec("UPDATE patients SET password = ? WHERE patient_id = ?", string(newHash), patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "密碼更新失敗: " + err.Error(),
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "密碼已成功更新",
	})
}
