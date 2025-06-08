package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBonusRoutes(r *gin.Engine) {
	b := r.Group("/api/v1/bonus")
	{
		// 報到
		b.PATCH("/appointments/:appointment_id/checkin", controllers.CheckInAppointment)
		// 請假與替代
		b.POST("/doctor_leaves", controllers.CreateLeave)
		b.GET("/doctor_leaves/:leave_id/substitutes", controllers.GetSubstitutes)
		// 即時佇列
		b.GET("/doctors/:doctor_id/queue", controllers.GetLiveQueue)
	}
}
