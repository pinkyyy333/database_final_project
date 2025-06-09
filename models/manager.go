package models

import "time"

type Manager struct {
	ManagerID uint      `gorm:"primaryKey;column:manager_id" json:"manager_id"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Manager) TableName() string {
	return "managers"
}
