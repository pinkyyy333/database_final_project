package controllers

import (
	"net/http"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

type AppointmentRequest struct {
	DepartmentID    uint32    `json:"department_id"`
	DoctorID        uint32    `json:"doctor_id"`
	PatientID       string    `json:"patient_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	ServiceType     string    `json:"service_type"`
}

func CreateAppointment(c *gin.Context) {
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if req.AppointmentTime.Before(time.Now().UTC()) {
		RespondError(c, http.StatusBadRequest, "預約時間必須在未來")
		return
	}
	// 衝突檢查
	var cnt int64
	if err := db.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND appointment_time = ?", req.DoctorID, req.AppointmentTime).
		Count(&cnt).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "查詢衝突失敗")
		return
	}
	if cnt > 0 {
		RespondError(c, http.StatusConflict, "時段已被預約")
		return
	}
	a := models.Appointment{
		DepartmentID:    req.DepartmentID,
		DoctorID:        req.DoctorID,
		PatientID:       req.PatientID,
		AppointmentTime: req.AppointmentTime,
		Status:          "booked",
		ServiceType:     req.ServiceType,
	}
	if err := db.DB.Create(&a).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "建立預約失敗")
		return
	}
	RespondCreated(c, gin.H{"appointment": a})
}

func GetPatientAppointments(c *gin.Context) {
	pid := c.Param("patient_id")
	var list []models.Appointment
	if err := db.DB.Where("patient_id = ?", pid).Find(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取得預約失敗")
		return
	}
	RespondOK(c, gin.H{"appointments": list})
}

func GetDoctorAppointments(c *gin.Context) {
	did := c.Param("doctor_id")
	var list []models.Appointment
	if err := db.DB.Where("doctor_id = ?", did).Find(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取得預約失敗")
		return
	}
	RespondOK(c, gin.H{"appointments": list})
}

type SlotInfo struct {
	Slot     string `json:"slot"`     // ISO8601 時間字串
	Count    int    `json:"count"`    // 已被預約人數
	Capacity int    `json:"capacity"` // 該時段最大可約人數 (slot_limit)
}

func GetAvailableSlots(c *gin.Context) {
	// 1. 解析路徑與查詢參數
	doctorID := c.Param("doctor_id")
	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "date 格式錯誤，請使用 YYYY-MM-DD",
		})
		return
	}

	// 2. 從 schedule_slots 表抓出該醫師該日所有排班時段
	var scheds []models.ScheduleSlot
	if err := db.DB.
		Where("doctor_id = ? AND date = ?", doctorID, dateStr).
		Find(&scheds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "讀取排班資料失敗",
		})
		return
	}

	// 3. 從 appointments 表統計每個時段的已被預約人數
	type CountRow struct {
		SlotTime string `gorm:"column:slot_time"` // TIME(appointment_time) 回傳
		Count    int    `gorm:"column:count"`
	}
	var counts []CountRow
	if err := db.DB.Table("appointments").
		Select("TIME(appointment_time) AS slot_time, COUNT(*) AS count").
		Where("doctor_id = ? AND DATE(appointment_time) = ?", doctorID, dateStr).
		Group("slot_time").
		Scan(&counts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "讀取預約統計失敗",
		})
		return
	}
	// 3.1 把統計結果放到 map 方便 lookup
	booked := make(map[string]int, len(counts))
	for _, r := range counts {
		booked[r.SlotTime] = r.Count
	}

	// 4. 組合最終要回傳的 infos
	infos := make([]SlotInfo, 0, len(scheds))
	for _, s := range scheds {
		// 4.1 先把字串 "HH:MM:SS" 解析成 time.Time
		parsed, err := time.Parse("15:04:05", s.SlotTime)
		if err != nil {
			continue
		}
		// 4.2 把它合併到當天的完整 time.Time
		t := time.Date(
			date.Year(), date.Month(), date.Day(),
			parsed.Hour(), parsed.Minute(), parsed.Second(),
			0, time.Local,
		)
		infos = append(infos, SlotInfo{
			Slot:     t.Format(time.RFC3339),
			Count:    booked[s.SlotTime],
			Capacity: s.SlotLimit,
		})
	}

	// 5. 回傳 JSON
	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"available_slots": infos,
	})
}

func UpdateAppointmentStatus(c *gin.Context) {
	id := c.Param("appointment_id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "參數錯誤")
		return
	}
	var a models.Appointment
	if err := db.DB.First(&a, id).Error; err != nil {
		RespondError(c, http.StatusNotFound, "找不到該預約")
		return
	}
	a.Status = req.Status
	if err := db.DB.Save(&a).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "更新狀態失敗")
		return
	}
	RespondOK(c, gin.H{"appointment": a})
}

func CancelAppointment(c *gin.Context) {
	id := c.Param("appointment_id")
	var a models.Appointment
	if err := db.DB.First(&a, id).Error; err != nil {
		RespondError(c, http.StatusNotFound, "找不到該預約")
		return
	}
	a.Status = "cancelled"
	if err := db.DB.Save(&a).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取消預約失敗")
		return
	}
	RespondOK(c, gin.H{"message": "預約已取消"})
}

func CheckInAppointment(c *gin.Context) {
	id := c.Param("appointment_id")
	var appt models.Appointment
	if err := db.DB.First(&appt, id).Error; err != nil {
		RespondError(c, http.StatusNotFound, "找不到該預約")
		return
	}
	now := time.Now().UTC()
	appt.CheckInTime = &now
	appt.Status = "checked_in"
	if err := db.DB.Save(&appt).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "報到失敗")
		return
	}
	RespondOK(c, gin.H{"appointment": appt})
}

func GetAvailableDoctors(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		RespondError(c, http.StatusBadRequest, "請提供 date 參數，格式 YYYY-MM-DD")
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "date 格式錯誤，需 YYYY-MM-DD")
		return
	}

	// 1) 從 schedule_slots 表撈出該 date 所有的排班時段
	var slots []models.ScheduleSlot
	if err := db.DB.
		Where("date = ?", date).
		Find(&slots).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取排班時段失敗")
		return
	}

	// 2) 從 slots 裡蒐集唯一的 doctor_id
	idMap := make(map[uint32]bool)
	for _, s := range slots {
		idMap[s.DoctorID] = true
	}
	ids := make([]uint32, 0, len(idMap))
	for id := range idMap {
		ids = append(ids, id)
	}

	// 3) 用抓到的 ids 去 doctors 表拿醫師資料
	var doctors []models.Doctor
	if err := db.DB.
		Where("doctor_id IN ?", ids).
		Find(&doctors).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取醫師資料失敗")
		return
	}

	// 4) 回傳 JSON
	RespondOK(c, gin.H{
		"success": true,
		"doctors": doctors,
	})
}

func GetAllAppointments(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		RespondError(c, http.StatusBadRequest, "請提供 date 參數")
		return
	}
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		RespondError(c, http.StatusBadRequest, "date 格式錯誤，需 YYYY-MM-DD")
		return
	}

	// ① 在函式內定義回傳用的 struct
	type AppointmentInfo struct {
		AppointmentID   uint32    `json:"appointment_id"`
		PatientName     string    `json:"patient_name"`
		DoctorName      string    `json:"doctor_name"`
		DepartmentName  string    `json:"department_name"`
		AppointmentTime time.Time `json:"appointment_time"`
		Status          string    `json:"status"`
	}

	// ② 一開始就初始化為長度為 0 的 slice
	list := make([]AppointmentInfo, 0)

	// ③ 用 GORM Scan 填值到 list
	if err := db.DB.
		Table("appointments AS a").
		Select(`
            a.appointment_id,
            p.patient_name     AS patient_name,
            d.doctor_name      AS doctor_name,
            dept.dept_name     AS department_name,
            a.appointment_time,
            a.status
        `).
		Joins("JOIN patients    AS p    ON p.patient_id    = a.patient_id").
		Joins("JOIN doctors     AS d    ON d.doctor_id     = a.doctor_id").
		Joins("JOIN departments AS dept ON dept.dept_id     = a.department_id").
		Where("DATE(a.appointment_time) = ?", dateStr).
		Scan(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取預約清單失敗")
		return
	}

	// ④ 回傳 JSON；即便 list 為空，也會是 [] 而不是 null
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"appointments": list,
	})
}
