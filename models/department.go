package models

import (
	"gorm.io/gorm"
)

// Department 代表科別
type Department struct {
	gorm.Model
	DeptName        string   `gorm:"column:dept_name"`
	DeptDescription string   `gorm:"column:dept_description"`
	Doctors         []Doctor `gorm:"foreignKey:DeptID"`
}
