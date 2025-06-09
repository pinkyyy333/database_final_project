package models

import "time"

// Schedule 代表每天某位醫師的排班
type Schedule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DoctorID  uint32    `gorm:"column:doctor_id" json:"doctor_id"`
	Date      time.Time `gorm:"column:date" json:"date"`
	SlotLimit int       `gorm:"column:slot_limit" json:"slot_limit"`
}

// TableName 明確指定對應到 schedules 表
func (Schedule) TableName() string {
	return "schedules"
}
