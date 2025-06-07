package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterManagerRoutes(r *gin.Engine) {
	m := r.Group("/api/v1/manager")
	{
		m.POST("/doctors", controllers.CreateDoctor)
		m.GET("/doctors", controllers.GetAllDoctors)
		m.PUT("/doctors/:id", controllers.UpdateDoctor)
		m.DELETE("/doctors/:id", controllers.DeleteDoctor)
		// 時段設定
		m.GET("/appointment_slots", controllers.GetAppointmentSlots)
		m.PUT("/appointment_slots", controllers.UpdateAppointmentSlots)
		// 管理員預約清單查詢
		m.GET("/appointments", controllers.GetAllAppointments)
		// 修改預約狀態（確認／取消／標記未出現）
		m.PATCH("/appointments/:id", controllers.UpdateAppointmentStatus)
		// 刪除預約
		m.DELETE("/appointments/:id", controllers.DeleteAdminAppointment)
		m.POST("/substitute",         controllers.AssignSubstitute)
		m.GET("/reports",             controllers.GenerateReport)
	}
}
