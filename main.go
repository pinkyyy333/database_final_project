// main.go
package controllers

import (
	"clinic-backend/db"
	"clinic-backend/middleware"
	"clinic-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	routes.RegisterDoctorRoutes(r)
	routes.RegisterAppointmentRoutes(r)


func main() {
  db.Init()                   // ← 這裡把 db.DB 設為 *sql.DB
  r := gin.Default()
  routes.RegisterDoctorRoutes(r)
  routes.RegisterAppointmentRoutes(r)
  // … 其他 middleware / route …
  r.Run(":8080")
}
