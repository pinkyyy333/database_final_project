package db

import (
	"log"
	"os"

	"clinic-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 建立 GORM 連線並自動 migrate 所有 models
func InitDB() {
	// 載入 .env
	if err := godotenv.Load(); err != nil {
		log.Println("[警告] 無法載入 .env，請確認環境變數已手動設定")
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

	if err := DB.AutoMigrate(
		&models.Patient{},
		//&models.Doctor{},
		//&models.Department{},
		//&models.Appointment{},
		//&models.Feedback{},
		//&models.Manager{},
	); err != nil {
		log.Fatalf("AutoMigrate 失敗: %v", err)
	}
}
