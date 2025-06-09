// controllers/slot_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/slots?doctor_id=123&month=2025-06
func GetAppointmentSlots(c *gin.Context) {
	doctorID := c.Query("doctor_id")
	month := c.Query("month") // 格式 "YYYY-MM"
	if doctorID == "" || month == "" {
		RespondError(c, http.StatusBadRequest, "doctor_id 與 month 為必填參數")
		return
	}

	// 1. 計算該月的起訖日
	const layout = "2006-01-02"
	start, _ := time.Parse(layout, month+"-01")
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// 2. 從 schedule_slots 撈出所有該醫師該月的記錄
	var slots []models.ScheduleSlot
	if err := db.DB.
		Where("doctor_id = ? AND date BETWEEN ? AND ?", doctorID, start, end).
		Find(&slots).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "讀取排班時段失敗")
		return
	}

	// 3. 依 date 分組，組成 sessions 字串陣列
	grouped := make(map[string][]string)
	for _, s := range slots {
		// 假設 s.SlotTime 格式為 "15:04:05"
		tm, _ := time.Parse("15:04:05", s.SlotTime)
		dateKey := s.Date.Format("2006-01-02")
		if tm.Hour() < 12 {
			grouped[dateKey] = append(grouped[dateKey], "morning")
		} else {
			grouped[dateKey] = append(grouped[dateKey], "afternoon")
		}
	}

	// 4. 將 map 轉成 []gin.H 回前端
	out := make([]gin.H, 0, len(grouped))
	for date, sess := range grouped {
		out = append(out, gin.H{
			"date":     date,
			"sessions": sess,
		})
	}

	RespondOK(c, gin.H{
		"doctor_id": doctorID,
		"month":     month,
		"slots":     out,
	})
}

// PUT /api/v1/slots
func UpdateAppointmentSlots(c *gin.Context) {
	const defaultLimit = 5

	// 1. Bind JSON payload
	var payload struct {
		DoctorID string                `json:"doctor_id"`
		Slots    []map[string][]string `json:"slots"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		RespondError(c, http.StatusBadRequest, "參數錯誤")
		return
	}

	// 2. 將 doctor_id 由字串轉 int
	idInt, err := strconv.Atoi(payload.DoctorID)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "doctor_id 格式錯誤")
		return
	}

	// 3. 刪除該醫師該月的舊資料
	//    取第一筆 slots map 的 key 裡的年月
	if len(payload.Slots) > 0 {
		for dateStr := range payload.Slots[0] {
			monthPrefix := dateStr[:7] + "-%" // e.g. "2025-06-%"
			db.DB.
				Where("doctor_id = ? AND date LIKE ?", idInt, monthPrefix).
				Delete(&models.ScheduleSlot{})
			break
		}
	}

	// 4. 插入新的排班紀錄
	for _, dayMap := range payload.Slots {
		for dateStr, sessions := range dayMap {
			// 4.1 解析日期字串
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}
			// 4.2 依每個 session 建立資料
			for _, sh := range sessions {
				var tm time.Time
				if sh == "morning" {
					tm, _ = time.Parse("15:04:05", "09:00:00")
				} else {
					tm, _ = time.Parse("15:04:05", "13:00:00")
				}

				rec := models.ScheduleSlot{
					DoctorID:  uint32(idInt),
					Date:      parsedDate,
					SlotTime:  tm.Format("15:04:05"),
					SlotLimit: defaultLimit,
				}
				db.DB.Create(&rec)
			}
		}
	}

	RespondOK(c, gin.H{"message": "排班已更新"})
}
