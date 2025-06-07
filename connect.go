// db/connection.go
package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// 讀取 .env 或環境變數
	// 範例 DSN: "user:pass@tcp(127.0.0.1:3306)/clinicdb?charset=utf8mb4&parseTime=True&loc=Local"
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

	// 如果你之後要自動建立/更新資料表，可以在這裡加 AutoMigrate
	// 例如：
	// DB.AutoMigrate(
	//   &models.Doctor{},
	//   &models.Patient{},
	//   &models.Appointment{},
	//   &models.Manager{},
	//   &models.Feedback{},
	// )
}
