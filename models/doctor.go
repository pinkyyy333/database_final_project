package models

import (
	"time"
)

// Doctor 對應資料庫 doctors table
type Doctor struct {
	DoctorID   uint32    `gorm:"column:doctor_id;primaryKey;autoIncrement:false" json:"doctor_id"`
	DeptID     uint32    `gorm:"column:dept_id" json:"dept_id"`
	DoctorName string    `gorm:"column:doctor_name" json:"doctor_name"`
	DoctorInfo string    `gorm:"column:doctor_info" json:"doctor_info"`
	Password   string    `gorm:"column:password" json:"-"` // 加密後密碼，不回傳給前端
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Gender     string    `gorm:"column:gender" json:"gender"`
	Edu        string    `gorm:"column:edu"    json:"edu"`
	HireDate   time.Time `gorm:"column:hire_date" json:"hire_date"`
	Phone      string    `gorm:"column:phone"  json:"phone"`
}

// TableName 指定資料表名稱
func (Doctor) TableName() string {
	return "doctors"
}
