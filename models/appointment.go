package models

import (
	"time"
)

// Appointment represents a booking or service
type Appointment struct {
    AppointmentID   int        `gorm:"primaryKey;column:Appointment_ID;autoIncrement" json:"appointment_id"`
    DeptID          int `gorm:"column:Dept_ID;type:int(11)" json:"dept_id"`
    DoctorID        string     `gorm:"column:Doctor_ID;type:char(5)" json:"doctor_id"`
    PatientID       string     `gorm:"column:Patient_ID;type:char(10)" json:"patient_id"`
    AppointmentTime time.Time  `gorm:"column:Appointment_Time" json:"appointment_time"`
    Status          string     `gorm:"column:Status;type:enum('booked','completed','cancelled','no_show');default:'booked'" json:"status"`
    
    ServiceType     string     `gorm:"column:Service_Type" json:"service_type"`           // 需確認資料庫有無此欄位
    CheckInTime     *time.Time `gorm:"column:Checkin_Time" json:"checkin_time,omitempty"` // 同上
}
