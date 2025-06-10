package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAppointmentRoutes(r *gin.Engine) {
	a := r.Group("/api/v1/appointments")
	{
		a.POST("", controllers.CreateAppointment)
		a.GET("/patient/:patient_id", controllers.GetPatientAppointments)
		a.GET("/doctor/:doctor_id", controllers.GetDoctorAppointments)
		a.GET("/doctor/:doctor_id/available", controllers.GetAvailableSlots)
		a.GET("/available-doctors", controllers.GetAvailableDoctors)
		a.PATCH("/:appointment_id/status", controllers.UpdateAppointmentStatus)
		a.DELETE("/:appointment_id", controllers.CancelAppointment)
		a.POST("/:appointment_id/checkin", controllers.CheckInAppointment)

		// 週班表相關
		a.GET("/schedule/months", controllers.GetScheduleMonths)
		a.GET("/schedule/weeks", controllers.GetScheduleWeeks)
		a.GET("/schedule", controllers.GetScheduleByWeek)
	}
}
