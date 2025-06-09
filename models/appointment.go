package models

import (
	"time"
)

// Appointment represents a booking or service
type Appointment struct {
	AppointmentID   uint       `gorm:"primaryKey;column:appointment_id" json:"appointment_id"`
	DepartmentID    uint       `gorm:"column:department_id" json:"department_id"`
	DoctorID        uint       `json:"doctor_id"`
	PatientID       string     `json:"patient_id"`
	AppointmentTime time.Time  `json:"appointment_time"`
	Status          string     `json:"status"`
	ServiceType     string     `gorm:"column:service_type" json:"service_type"`           // 新增：服務類型（consult/vaccine/...）
	CheckInTime     *time.Time `gorm:"column:checkin_time" json:"checkin_time,omitempty"` // 新增：報到時間
}
