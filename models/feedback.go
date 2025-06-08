// models/feedback.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Feedback 代表病患對一次預約的評價
type Feedback struct {
	gorm.Model
	AppointmentID  uint      `json:"appointment_id"`
	FeedbackRating int       `json:"feedback_rating"`
	PatientComment string    `json:"patient_comment"`
	CreatedAt      time.Time `json:"created_at"`

	// 關聯
	Appointment Appointment `gorm:"foreignKey:AppointmentID"`
}
