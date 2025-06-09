// routes/slot_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSlotRoutes(r *gin.Engine) {
	s := r.Group("/api/v1/slots")
	{
		s.GET("", controllers.GetAppointmentSlots)    // 查詢 doctor_id & month 的可用時段
		s.PUT("", controllers.UpdateAppointmentSlots) // 更新 doctor_id 的時段
	}
}
