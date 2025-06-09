package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterDoctorRoutes(r *gin.Engine) {
	d := r.Group("/api/v1/doctors")
	{
		d.POST("", controllers.CreateDoctor)
		d.POST("/login", controllers.LoginDoctor)
		d.GET("/:id/schedule", controllers.GetDoctorSchedule)
		d.PUT("/appointments/:id/status", controllers.UpdateAppointmentStatus)
		d.GET("/:id/patients/:patient_id/records", controllers.GetPatientRecords)
	}
}
