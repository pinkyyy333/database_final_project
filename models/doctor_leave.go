package models

import "time"

// DoctorLeave 代表醫師請假與替代安排
type DoctorLeave struct {
	LeaveID            uint      `gorm:"primaryKey;column:leave_id" json:"leave_id"`
	DoctorID           uint      `json:"doctor_id"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	SubstituteDoctorID *uint     `json:"substitute_doctor_id"` // 替代醫師
}
