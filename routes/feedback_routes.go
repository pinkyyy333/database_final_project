// routes/feedback_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFeedbackRoutes(r *gin.Engine) {
	// feedbacks endpoint
	f := r.Group("/api/v1/feedbacks")
	{
		f.POST("", controllers.CreateFeedback) // POST  /api/v1/feedbacks
	}
	// doctor feedbacks
	d := r.Group("/api/v1/doctors")
	{
		d.GET("/:id/feedbacks", controllers.GetDoctorFeedbacks)
	}
}
