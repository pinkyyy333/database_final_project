package models

type Manager struct {
	ManagerID uint   `gorm:"primaryKey;column:manager_id" json:"manager_id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
}
