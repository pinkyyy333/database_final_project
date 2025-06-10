package models

import "time"

type ScheduleSlot struct {
	SlotID    uint      `gorm:"primaryKey;column:slot_id" json:"slot_id"`
	DoctorID  uint32    `gorm:"column:doctor_id"    json:"doctor_id"`
	Date      time.Time `gorm:"column:date"         json:"date"`
	SlotTime  string    `gorm:"column:slot_time"    json:"time"`
	SlotLimit int64     `gorm:"column:slot_limit"   json:"-"`
	DeptID    uint32    `gorm:"column:dept_id" json:"dept_id"`
}

func (ScheduleSlot) TableName() string {
	return "schedule_slots"
}
