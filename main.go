package main

import (
	"appointment-system/config"
	"appointment-system/models"
	"appointment-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	config.ConnectDatabase()

	// Run migrations
	config.DB.AutoMigrate(&models.Patient{}, &models.Doctor{}, &models.Appointment{})

	// Initialize Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	// Start the server
	r.Run(":8080")
}
