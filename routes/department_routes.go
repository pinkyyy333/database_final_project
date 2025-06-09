// routes/department_routes.go
package routes

import (
	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterDepartmentRoutes 提供列出科別，以及依科別查醫師
func RegisterDepartmentRoutes(r *gin.Engine) {
	d := r.Group("/api/v1/departments")
	{
		// 取得所有科別
		d.GET("", controllers.GetAllDepartments)
		// 取得某科別底下的醫師清單
		d.GET("/:dept_id/doctors", controllers.GetDoctorsByDepartment)
	}
}
