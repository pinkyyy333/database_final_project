package routes

import (
	"clinic-backend/controllers"
	"clinic-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(r *gin.Engine) {
	p := r.Group("/api/v1/patients")
	{
		p.POST("/register", controllers.RegisterPatient)
		p.POST("/login", controllers.LoginPatient)

		p.Use(middleware.AuthMiddleware())
		p.GET("/profile", controllers.GetPatientProfile)
		p.PUT("/profile", controllers.UpdatePatientProfile)
		p.PUT("/password", controllers.ChangePatientPassword)
	}
}
