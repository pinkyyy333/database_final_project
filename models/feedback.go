// models/feedback.go
package models

import (
	"time"
)

// Feedback 代表病患對一次預約的評價
type Feedback struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	AppointmentID  uint      `json:"appointment_id"`
	FeedbackRating int       `json:"feedback_rating"`
	PatientComment string    `json:"patient_comment"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
}
