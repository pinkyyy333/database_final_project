package services

import (
	"clinic-backend/db"
	"clinic-backend/models"
	"time"

	"github.com/robfig/cron/v3"
)

// StartNoShowCron 啟動排程，每 15 分鐘檢測未出現
func StartNoShowCron() {
	c := cron.New()
	c.AddFunc("*/15 * * * *", detectNoShows)
	c.Start()
}

func detectNoShows() {
	now := time.Now().UTC()
	var list []models.Appointment
	db.DB.Where("status = ? AND appointment_time < ?", "booked", now).
		Find(&list)
	for _, appt := range list {
		// 若尚未報到（CheckInTime 為 nil），標記 no_show
		if appt.CheckInTime == nil {
			appt.Status = "no_show"
			db.DB.Save(&appt)
		}
	}
}
