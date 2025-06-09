// routes/leave_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterLeaveRoutes(r *gin.Engine) {
	l := r.Group("/api/v1/leaves")
	{
		l.POST("", controllers.CreateLeave)                         // 建立請假
		l.GET("/:leave_id/substitutes", controllers.GetSubstitutes) // 查詢科別替代醫師
	}
}
