// routes/appointment_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAppointmentRoutes 把新的「取消」與「可用時段」也一併掛上
func RegisterAppointmentRoutes(r *gin.Engine) {
	a := r.Group("/api/v1/appointments")
	{
		a.POST("", controllers.CreateAppointment)
		a.GET("/patient/:patient_id", controllers.GetPatientAppointments)
		a.GET("/doctor/:doctor_id", controllers.GetDoctorAppointments)
		// 新增：查詢醫師在某日的可預約時段（query param: ?date=2025-06-10）
		a.GET("/doctor/:doctor_id/available", controllers.GetAvailableSlots)
		a.PATCH("/:appointment_id/status", controllers.UpdateAppointmentStatus)
		// 新增：取消預約（將狀態標為 "cancelled"）
		a.DELETE("/:appointment_id", controllers.CancelAppointment)
	}
}
