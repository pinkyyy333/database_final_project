package models

import "time"

// VaccineAppointment 用來記錄疫苗接種預約資料
type VaccineAppointment struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Birthdate   time.Time `gorm:"type:date;not null" json:"birthdate"`
	HealthCard  string    `gorm:"not null;unique" json:"healthcard"`
	IDNumber    string    `gorm:"not null" json:"idnumber"`
	VaccineDate time.Time `gorm:"type:date;not null" json:"vaccine_date"`
	VaccineType string    `gorm:"not null" json:"vaccine_type"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (VaccineAppointment) TableName() string {
	return "vaccine_appointments"
}
