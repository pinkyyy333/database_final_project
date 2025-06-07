package services

import (
	"fmt"
	"time"

	"clinic-backend/db"
	"clinic-backend/models"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

// StartReminderCron 啟動排程，每 30 分鐘發送即將到來預約提醒
func StartReminderCron() {
	c := cron.New()
	c.AddFunc("0,30 * * * *", sendReminders)
	c.Start()
}

func sendReminders() {
	// 查詢 1 小時內的預約
	horizon := time.Now().UTC().Add(1 * time.Hour)
	var list []models.Appointment
	db.DB.Where("status = ? AND appointment_time BETWEEN ? AND ?", "booked", time.Now().UTC(), horizon).
		Find(&list)
	for _, appt := range list {
		// 假設 models.Patient 有 Email 欄位
		var patient models.Patient
		db.DB.First(&patient, "patient_id = ?", appt.PatientID)
		// 發信
		msg := gomail.NewMessage()
		msg.SetHeader("From", "noreply@clinic.com")
		msg.SetHeader("To", patient.PatientPhone+"@sms.gateway.example.com") // 或 Email
		msg.SetHeader("Subject", "預約提醒")
		body := fmt.Sprintf("您在 %s 的預約即將開始（%s）", appt.AppointmentTime.Format(time.RFC3339), appt.ServiceType)
		msg.SetBody("text/plain", body)
		dialer := gomail.NewDialer("smtp.example.com", 587, "user", "pass")
		dialer.DialAndSend(msg)
	}
}
