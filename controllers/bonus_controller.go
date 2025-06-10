package controllers

import (
	"clinic-backend/db"
	"clinic-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/bonus/vaccine/quota?date=YYYY-MM-DD
func GetVaccineQuota(c *gin.Context) {
	date := c.Query("date")
	count := db.GetVaccineCountByDate(date)
	c.JSON(http.StatusOK, gin.H{
		"date":      date,
		"booked":    count,
		"remaining": 20 - count,
	})
}

// POST /api/v1/bonus/vaccine
func CreateVaccineAppointment(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Birthdate   string `json:"birthdate"` // "2006-01-02"
		HealthCard  string `json:"healthcard"`
		IDNumber    string `json:"idnumber"`
		VaccineDate string `json:"vaccine_date"` // "2006-01-02"
		VaccineType string `json:"vaccine_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error(), "code": 400})
		return
	}

	// 解析日期
	bd, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "出生日期格式錯誤", "code": 400})
		return
	}
	vd, err := time.Parse("2006-01-02", req.VaccineDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "接種日期格式錯誤", "code": 400})
		return
	}

	// 1) 額度檢查
	booked := db.GetVaccineCountByDate(req.VaccineDate)
	if booked >= 20 {
		c.JSON(http.StatusConflict, gin.H{"error": true, "message": "該日已滿額，請選擇其他日期", "code": 409})
		return
	}

	// 2) 重複預約檢查：同一身分證當日不可重複
	var dupCount int64
	if err := db.DB.
		Model(&models.VaccineAppointment{}).
		Where("id_number = ? AND vaccine_date = ?", req.IDNumber, vd).
		Count(&dupCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "伺服器內部錯誤", "code": 500})
		return
	}
	if dupCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": true, "message": "您已於此日期預約過疫苗", "code": 409})
		return
	}

	// 建立資料
	appointment := models.VaccineAppointment{
		Name:        req.Name,
		Birthdate:   bd,
		HealthCard:  req.HealthCard,
		IDNumber:    req.IDNumber,
		VaccineDate: vd,
		VaccineType: req.VaccineType,
	}
	if err := db.CreateVaccineAppointment(appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "預約失敗，請稍後重試", "code": 500})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "預約成功", "code": 201})
}
