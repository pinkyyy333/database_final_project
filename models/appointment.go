package models

import (
	"time"
)

// Appointment represents a booking or service
type Appointment struct {
	AppointmentID   int       `gorm:"column:Appointment_ID;primaryKey" json:"appointment_id"`
	DepartmentID    int       `gorm:"column:Dept_ID" json:"department_id"`
	DoctorID        string    `gorm:"column:Doctor_ID" json:"doctor_id"`
	PatientID       string    `gorm:"column:Patient_ID" json:"patient_id"`
	AppointmentTime time.Time `gorm:"column:Appointment_Time" json:"appointment_time"`
	Status          string    `gorm:"column:Status" json:"status"`
	ServiceType     string     `gorm:"column:service_type" json:"service_type"`           // 新增：服務類型（consult/vaccine/...）
	CheckInTime     *time.Time `gorm:"column:checkin_time" json:"checkin_time,omitempty"` // 新增：報到時間
}

// TableName overrides the default table name
func (Appointment) TableName() string {
	return "Appointments"
}