package handlers

import (
	"appointment-system/config"
	"appointment-system/models"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func CreateAppointment(c *gin.Context) {
	var input models.Appointment
	// var input struct {
	// 	PatientEmail   string `json:"patient_email" binding:"required,email"`
	// 	DoctorEmail    string `json:"doctor_email" binding:"required,email"`
	// 	Specialization string `json:"specialization" binding:"required"`
	// 	TimeSlot       struct {
	// 		StartTime string `json:"start_time" binding:"required"`
	// 		EndTime   string `json:"end_time" binding:"required"`
	// 	} `json:"timeSlot" binding:"required"`
	// }

	// Log the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read request body."})
		return
	}
	fmt.Println("Raw Request Body:", string(body))

	// Reset the request body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Invalid data to create", "error": err.Error()})
		return
	}

	// Check if patient exists
	var patient models.Patient
	if err := config.DB.Where("email = ?", input.PatientEmail).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found."})
		return
	}

	// Check if doctor exists
	var doctor models.Doctor
	if err := config.DB.Where("email = ?", input.DoctorEmail).First(&doctor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found."})
		return
	}

	// Check for existing appointment
	var existingAppointment models.Appointment
	if config.DB.Where("patient_id = ? AND status <> ? AND doctor_id = ?", patient.ID, "CANCELLED", doctor.ID).First(&existingAppointment).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Patient already has an appointment booked with this doctor."})
		return
	}

	// Log the start and end times to check their values
	fmt.Println("Start Time:", input.TimeSlot.StartTime)
	fmt.Println("End Time:", input.TimeSlot.EndTime)

	// Convert start and end times to minutes
	formattedStartTime, err := timeToMinutes(input.TimeSlot.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	formattedEndTime, err := timeToMinutes(input.TimeSlot.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if end time is greater than start time
	if formattedEndTime <= formattedStartTime {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "End time must be greater than start time"})
		return
	}

	if formattedEndTime == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// Check for overlapping appointments
	var overlappingAppointment models.Appointment
	if config.DB.Where("doctor_id = ? AND ((timeStamp->'$.startTime' < ? AND timeStamp->'$.endTime' > ?) OR (timeStamp->'$.endTime' > ? AND timeStamp->'$.startTime' < ?))", doctor.ID, formattedEndTime, formattedStartTime, formattedStartTime, formattedEndTime).First(&overlappingAppointment).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Time slot is already booked."})
		return
	}

	// Create the appointment
	newAppointment := models.Appointment{
		PatientID:      patient.ID,
		DoctorID:       doctor.ID,
		Specialization: input.Specialization,
		Status:         "BOOKED",
		TimeSlot: models.TimeStamp{
			StartTime: cast.ToString(formattedStartTime),
			EndTime:   cast.ToString(formattedEndTime),
		},
		PatientEmail: patient.Email,
		DoctorEmail:  doctor.Email,
	}

	// Save the appointment to the database
	if err := config.DB.Create(&newAppointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Appointment created successfully.",
		"appointment": newAppointment,
	})
}

// timeToMinutes converts a time string in RFC3339 format to total minutes
func timeToMinutes(timeStr string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing time: %w", err)
	}
	return int64(t.Hour()*60 + t.Minute()), nil
}

// ViewAllAppointments retrieves all appointments for a doctor
func ViewAllAppointments(c *gin.Context) {
	doctorEmail := c.Query("doctor")
	var appointments []models.Appointment
	if err := config.DB.Where("doctor_email = ? AND status = ?", doctorEmail, "BOOKED").Find(&appointments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No appointments found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "Success",
		"model": gin.H{"appointments": appointments},
	})
}

// ViewAppointmentDetails retrieves details for a specific appointment
func ViewAppointmentDetails(c *gin.Context) {
	patientEmail := c.Query("patient_email")
	var appointment models.Appointment
	if err := config.DB.Where("patient_email = ?", patientEmail).First(&appointment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "Success",
		"model": gin.H{"appointmentDetails": appointment},
	})
}

// CancelAppointment cancels a specific appointment
func CancelAppointment(c *gin.Context) {
	var req struct {
		Patient  string           `json:"patient_email"`
		Doctor   string           `json:"doctor_email"`
		TimeSlot models.TimeStamp `json:"timeSlot"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var appointment models.Appointment
	if err := config.DB.Where("patient_email = ? AND doctor_email = ? AND status = ?", req.Patient, req.Doctor, "BOOKED").First(&appointment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	appointment.Status = "Cancelled"
	if err := config.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cancel appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Appointment Cancelled"})
}
