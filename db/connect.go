// db/connect.go
package db

import (
	"clinic-backend/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("[警告] 無法載入 .env")
	}

	dsn := os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASS") + "@tcp(" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME") +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("GORM 連線失敗: %v", err)
	}

	// Debug 模式：印出完整 SQL
	DB = DB.Debug()

	if err := DB.AutoMigrate(
		&models.Patient{},
		&models.Department{},
		&models.Doctor{},
		&models.Appointment{},
		&models.Feedback{},
		&models.Manager{},
		&models.VaccineAppointment{},
	); err != nil {
		log.Fatalf("AutoMigrate 失敗: %v", err)
	}
}

// （已移除 GetVaccineCountByDate 及 CreateVaccineAppointment）
