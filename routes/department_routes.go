// routes/department_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterDepartmentRoutes(r *gin.Engine) {
	d := r.Group("/api/v1/departments")
	{
		d.GET("", controllers.GetAllDepartments)
		d.GET("/:dept_id/doctors", controllers.GetDoctorsByDepartment)
	}
}
