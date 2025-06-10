package controllers

import (
	"fmt"
	"net/http"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/gin-gonic/gin"
)

// GetScheduleMonths 回傳可選月份清單
func GetScheduleMonths(c *gin.Context) {
	var months []string
	db.DB.
		Model(&models.Schedule{}).
		Distinct("DATE_FORMAT(date, '%Y-%m')").
		Pluck("DATE_FORMAT(date, '%Y-%m')", &months)
	c.JSON(http.StatusOK, gin.H{"months": months})
}

// GetScheduleWeeks 回傳該月各週起始日列表
// GET /api/v1/appointments/schedule/weeks?month=2025-06
func GetScheduleWeeks(c *gin.Context) {
	month := c.Query("month")
	var weeks []string
	db.DB.
		Model(&models.Schedule{}).
		Distinct("DATE_FORMAT(date - INTERVAL (DAYOFWEEK(date)-1) DAY, '%Y-%m-%d')").
		Where("DATE_FORMAT(date, '%Y-%m') = ?", month).
		Order("1").
		Pluck("DATE_FORMAT(date - INTERVAL (DAYOFWEEK(date)-1) DAY, '%Y-%m-%d')", &weeks)
	c.JSON(http.StatusOK, gin.H{"weeks": weeks})
}

// GetScheduleByWeek 回傳指定週（起始日）整週的班表
// GET /api/v1/appointments/schedule?week=2025-06-01
func GetScheduleByWeek(c *gin.Context) {
	week := c.Query("week")
	start, _ := time.Parse("2006-01-02", week)
	end := start.AddDate(0, 0, 6)

	// 1) 撈出這週所有排班 slot
	var slots []models.ScheduleSlot
	db.DB.
		Where("date BETWEEN ? AND ?", start, end).
		Order("slot_time, dept_id").
		Find(&slots)

	// 2) 組 dates & days
	dates := make([]string, 7)
	days := []string{"日", "一", "二", "三", "四", "五", "六"}
	for i := 0; i < 7; i++ {
		d := start.AddDate(0, 0, i)
		dates[i] = d.Format("01/02")
	}

	// 3) 撈所有醫師對應 Dept（如要 Dept 名稱可 JOIN Department table）
	var docs []models.Doctor
	db.DB.Select("doctor_id, doctor_name, dept_id").Find(&docs)
	type DocInfo struct {
		Name  string `json:"name"`
		Full  bool   `json:"full"`
		Reg   string `json:"reg,omitempty"`
		Start string `json:"start,omitempty"`
	}
	docMap := map[uint32]DocInfo{}
	for _, d := range docs {
		docMap[d.DoctorID] = DocInfo{Name: d.DoctorName}
	}

	// 4) 組 rowMap
	type Row struct {
		Time     string    `json:"time"`
		DeptName string    `json:"dept_name"`
		Doctors  []DocInfo `json:"doctors"`
	}
	rowsMap := map[string]*Row{}
	for _, s := range slots {
		key := fmt.Sprintf("%s_%d", s.SlotTime, s.DeptID)
		if rowsMap[key] == nil {
			rowsMap[key] = &Row{
				Time:     s.SlotTime,
				DeptName: "", // 若有 Department table，可在此取得科別名稱
				Doctors:  make([]DocInfo, 7),
			}
		}
		idx := int(s.Date.Sub(start).Hours() / 24)

		// 計算已掛號人數
		var cnt int64
		db.DB.Model(&models.Appointment{}).
			Where("doctor_id = ? AND DATE(appointment_time)=? AND TIME(appointment_time)=?",
				s.DoctorID, s.Date, s.SlotTime).
			Count(&cnt)

		info := docMap[s.DoctorID]
		info.Full = cnt >= s.SlotLimit
		if cnt > 0 {
			info.Reg = fmt.Sprintf("已掛到%d號", cnt)
		}
		rowsMap[key].Doctors[idx] = info
	}

	// 5) 轉 slice 回傳
	resp := make([]Row, 0, len(rowsMap))
	for _, r := range rowsMap {
		resp = append(resp, *r)
	}

	c.JSON(http.StatusOK, gin.H{
		"dates":    dates,
		"days":     days,
		"schedule": resp,
	})
}
