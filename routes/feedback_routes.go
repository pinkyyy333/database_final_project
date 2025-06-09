// routes/feedback_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFeedbackRoutes(r *gin.Engine) {
	f := r.Group("/api/v1/feedbacks")
	{
		f.POST("", controllers.CreateFeedback)
	}
	// 供醫師取得評價
	r.GET("/api/v1/doctors/:doctor_id/feedbacks", controllers.GetDoctorFeedbacks)
}
