// routes/queue_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterQueueRoutes(r *gin.Engine) {
	q := r.Group("/api/v1")
	{
		q.GET("/doctors/:doctor_id/queue", controllers.GetLiveQueue)
	}
}
