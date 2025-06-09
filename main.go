package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"clinic-backend/db"
	"clinic-backend/routes"
	services "clinic-backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 靜態檔的資料夾
	publicDir := "./public"
	if _, err := os.Stat(publicDir); err != nil {
		panic("找不到 public 目錄: " + err.Error())
	}

	r := gin.Default()
	r.Use(cors.Default())
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system env")
	}
	// --- 先註冊所有 /api/... 路由 ---
	db.InitDB()
	go services.StartReminderCron()
	go services.StartNoShowCron()

	routes.RegisterAppointmentRoutes(r)
	routes.RegisterDoctorRoutes(r)
	routes.RegisterPatientRoutes(r)
	routes.RegisterManagerRoutes(r)
	routes.RegisterBonusRoutes(r)
	routes.RegisterFeedbackRoutes(r)
	routes.RegisterDepartmentRoutes(r)
	routes.RegisterSlotRoutes(r)

	// --- 再用 NoRoute 來處理所有非 /api/... 的請求，回傳對應的靜態檔 ---
	r.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		// 如果是 API 路徑，就直接 404
		if strings.HasPrefix(reqPath, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}

		// 轉成相對檔名，預設 "/" 走 index.html
		file := strings.TrimPrefix(reqPath, "/")
		if file == "" {
			file = "index.html"
		}

		fullPath := filepath.Join(publicDir, file)
		if _, err := os.Stat(fullPath); err != nil {
			// 檔案不存在，回 404
			c.Status(http.StatusNotFound)
			return
		}

		// 回傳靜態檔案
		c.File(fullPath)
	})

	// 啟動服務
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
