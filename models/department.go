package models

import "time"

// Department 代表科別
type Department struct {
	DeptID          uint      `gorm:"primaryKey;column:dept_id" json:"dept_id"`
	DeptName        string    `gorm:"column:dept_name" json:"dept_name"`
	DeptDescription string    `gorm:"column:dept_description" json:"dept_description"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
	Doctors         []Doctor  `gorm:"foreignKey:dept_id" json:"doctors"`
}
