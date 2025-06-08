package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterManagerRoutes(r *gin.Engine) {
	m := r.Group("/api/v1/manager")
	{
		m.POST("/login", controllers.LoginManager)
		m.POST("/doctors", controllers.CreateDoctor)
		m.GET("/doctors", controllers.GetAllDoctors)
		m.PUT("/doctors/:id", controllers.UpdateDoctor)
		m.DELETE("/doctors/:id", controllers.DeleteDoctor)
		m.GET("/reports", controllers.GenerateReport)
	}
}
