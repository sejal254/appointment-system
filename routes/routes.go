package routes

import (
	"appointment-system/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	//Appointment routes
	r.POST("/appointments", handlers.CreateAppointment)
	r.GET("/appointments", handlers.ViewAllAppointments)
	r.GET("/appointments/details", handlers.ViewAppointmentDetails)
	r.POST("/appointments/cancel", handlers.CancelAppointment)

	//user routes
	r.POST("/create/patient", handlers.AddPatient)
	r.POST("/create/doctor", handlers.AddDoctor)

}
