package handlers

import (
	"appointment-system/config"
	"appointment-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddPatient handles adding a new patient
func AddPatient(c *gin.Context) {
	var input models.Patient

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required fields."})
		return
	}

	// Check if patient already exists
	var patient models.Patient
	if err := config.DB.Where("email = ?", input.Email).First(&patient).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Patient Already Exists"})
		return
	}

	// Create a new patient instance
	newPatient := models.Patient{
		Name:  input.Name,
		Email: input.Email,
	}

	// Save the new patient to the database
	if err := config.DB.Create(&newPatient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add patient."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient added successfully!",
		"patient": newPatient,
	})
}

// AddDoctor handles adding a new doctor
func AddDoctor(c *gin.Context) {
	var input models.Doctor

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and specialization are required."})
		return
	}

	// Check if doctor already exists
	var doctor models.Doctor
	if err := config.DB.Where("email = ?", input.Email).First(&doctor).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Doctor Already Exists"})
		return
	}

	// Create a new doctor instance
	newDoctor := models.Doctor{
		Name:           input.Name,
		Email:          input.Email,
		Specialization: input.Specialization,
	}

	// Save the new doctor to the database
	if err := config.DB.Create(&newDoctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add doctor."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Doctor added successfully!",
		"doctor":  newDoctor,
	})
}
