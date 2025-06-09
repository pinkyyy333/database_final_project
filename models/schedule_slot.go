package models

import "time"

// ScheduleSlot 代表某位醫師在某日某時段的排班
type ScheduleSlot struct {
	ID        uint      `gorm:"primaryKey;column:id"     json:"id"`
	DoctorID  uint32    `gorm:"column:doctor_id"         json:"doctor_id"`
	Date      time.Time `gorm:"column:date"              json:"date"`
	SlotTime  string    `gorm:"column:slot_time"         json:"slot_time"`
	SlotLimit int       `gorm:"column:slot_limit"        json:"capacity"` // 回傳時用 capacity
}

// TableName 指定對應到 schedule_slots
func (ScheduleSlot) TableName() string {
	return "schedule_slots"
}
