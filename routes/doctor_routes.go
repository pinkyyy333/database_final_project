package routes

import (
	"clinic-backend/controllers"
	"clinic-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDoctorRoutes(r *gin.Engine) {
	d := r.Group("/api/v1/doctors", middleware.AuthMiddleware())
	{
		d.GET("/:id/schedule", controllers.GetDoctorSchedule)
		d.PUT("/appointments/:id/status", controllers.UpdateAppointmentStatus)
		d.GET("/:id/patients/:patient_id/records", controllers.GetPatientRecords)
	}
}
