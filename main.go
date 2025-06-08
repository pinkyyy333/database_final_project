package main

import (
	"clinic-backend/db"
	"clinic-backend/routes"
	services "clinic-backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	db.InitDB()

	// 啟動背景任務
	go services.StartReminderCron()
	go services.StartNoShowCron()

	// 註冊路由
	routes.RegisterAppointmentRoutes(r)
	routes.RegisterDoctorRoutes(r)
	routes.RegisterPatientRoutes(r)
	routes.RegisterManagerRoutes(r)
	routes.RegisterBonusRoutes(r)

	r.Run()
}
