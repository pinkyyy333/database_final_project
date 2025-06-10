package controllers

import (
	"net/http"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// AppointmentRequest 綁定前端的預約請求
type AppointmentRequest struct {
	DepartmentID    uint32    `json:"department_id"`
	DoctorID        uint32    `json:"doctor_id"`
	PatientID       string    `json:"patient_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	ServiceType     string    `json:"service_type"`
}

// CreateAppointment 建立新預約並檢查衝突
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

// GetPatientAppointments 取得某病患所有預約
func GetPatientAppointments(c *gin.Context) {
	pid := c.Param("patient_id")
	var list []models.Appointment
	if err := db.DB.Where("patient_id = ?", pid).Find(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取得預約失敗")
		return
	}
	RespondOK(c, gin.H{"appointments": list})
}

// GetDoctorAppointments 取得某醫師所有預約
func GetDoctorAppointments(c *gin.Context) {
	did := c.Param("doctor_id")
	var list []models.Appointment
	if err := db.DB.Where("doctor_id = ?", did).Find(&list).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "取得預約失敗")
		return
	}
	RespondOK(c, gin.H{"appointments": list})
}

// SlotInfo 用來回傳各時段的已預約數與容量
type SlotInfo struct {
	Slot     string `json:"slot"`     // ISO8601 時間字串
	Count    int    `json:"count"`    // 已被預約人數
	Capacity int64  `json:"capacity"` // 該時段最大可約人數
}

// GetAvailableSlots 查某醫師於某日的所有可預約時段
func GetAvailableSlots(c *gin.Context) {
	doctorID := c.Param("doctor_id")
	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "date 格式錯誤，請使用 YYYY-MM-DD")
		return
	}

	// 撈排班
	var scheds []models.ScheduleSlot
	if err := db.DB.
		Where("doctor_id = ? AND date = ?", doctorID, date).
		Find(&scheds).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取排班資料失敗")
		return
	}

	// 統計已預約數
	type countRow struct {
		SlotTime string `gorm:"column:slot_time"`
		Count    int    `gorm:"column:count"`
	}
	var counts []countRow
	if err := db.DB.Table("appointments").
		Select("TIME(appointment_time) AS slot_time, COUNT(*) AS count").
		Where("doctor_id = ? AND DATE(appointment_time) = ?", doctorID, dateStr).
		Group("slot_time").
		Scan(&counts).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取預約統計失敗")
		return
	}
	booked := make(map[string]int, len(counts))
	for _, r := range counts {
		booked[r.SlotTime] = r.Count
	}

	// 組回傳
	infos := make([]SlotInfo, 0, len(scheds))
	for _, s := range scheds {
		parsed, _ := time.Parse("15:04:05", s.SlotTime)
		t := time.Date(
			date.Year(), date.Month(), date.Day(),
			parsed.Hour(), parsed.Minute(), parsed.Second(),
			0, time.UTC,
		)
		infos = append(infos, SlotInfo{
			Slot:     t.Format(time.RFC3339),
			Count:    booked[s.SlotTime],
			Capacity: s.SlotLimit,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"available_slots": infos,
	})
}

// UpdateAppointmentStatus 更新某預約狀態
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

// CancelAppointment 將某預約標記為 cancelled
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

// CheckInAppointment 報到某預約，設定狀態為 checked_in 並記錄時間
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

// GetAvailableDoctors 查某日(YYYY-MM-DD)有排班的醫師列表
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
	var slots []models.ScheduleSlot
	if err := db.DB.Where("date = ?", date).Find(&slots).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取排班時段失敗")
		return
	}
	set := make(map[uint32]struct{})
	for _, s := range slots {
		set[s.DoctorID] = struct{}{}
	}
	ids := make([]uint32, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	var doctors []models.Doctor
	if err := db.DB.Where("doctor_id IN ?", ids).Find(&doctors).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取醫師資料失敗")
		return
	}
	RespondOK(c, gin.H{"doctors": doctors})
}

// 以下還有 GetAllAppointments、GetScheduleMonths、GetScheduleWeeks、GetScheduleByWeek 等函式…

// GetAllAppointments 列出某日所有預約(含病患/醫師/科別名稱)
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

	type AppointmentInfo struct {
		AppointmentID   uint32    `json:"appointment_id"`
		PatientName     string    `json:"patient_name"`
		DoctorName      string    `json:"doctor_name"`
		DepartmentName  string    `json:"department_name"`
		AppointmentTime time.Time `json:"appointment_time"`
		Status          string    `json:"status"`
	}
	list := make([]AppointmentInfo, 0)
	err := db.DB.Table("appointments AS a").
		Select("a.appointment_id, p.patient_name AS patient_name, d.doctor_name AS doctor_name, dept.dept_name AS department_name, a.appointment_time, a.status").
		Joins("JOIN patients AS p ON p.patient_id = a.patient_id").
		Joins("JOIN doctors AS d ON d.doctor_id = a.doctor_id").
		Joins("JOIN departments AS dept ON dept.dept_id = a.department_id").
		Where("DATE(a.appointment_time) = ?", dateStr).
		Scan(&list).Error
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取預約清單失敗")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"appointments": list,
	})
}
