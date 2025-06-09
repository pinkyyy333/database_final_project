package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterManagerRoutes(r *gin.Engine) {

	m := r.Group("/api/v1/manager")
	{
		// 管理員註冊
		// 管理員登入
		m.POST("/login", controllers.LoginManager)
		// （可選）查所有管理員
		m.GET("", controllers.GetAllManagers)
		// 更新管理員帳號或密碼
		m.PUT("/:id", controllers.UpdateManager)
		// 刪除管理員
		m.DELETE("/:id", controllers.DeleteManager)
		// 產生報表
		m.GET("/reports", controllers.GenerateReport)
		// ※ 全院預約清單（指定日期）
		m.GET("/appointments", controllers.GetAllAppointments)

		// ※ 更新某筆預約狀態（取消、標記未出現等）
		m.PATCH("/appointments/:appointment_id", controllers.UpdateAppointmentStatus)
	}
}
