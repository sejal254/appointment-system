package models

import "gorm.io/gorm"

type Appointment struct {
	gorm.Model
	PatientEmail   string    `json:"patient_email"`
	DoctorEmail    string    `json:"doctor_email"`
	PatientID      uint      `json:"patient_id" gorm:"not null"`         // Foreign key
	DoctorID       uint      `json:"doctor_id" gorm:"not null"`          // Foreign key
	TimeSlot       TimeStamp `json:"timeStamp" gorm:"embedded;not null"` // Embed TimeStamp struct
	Specialization string    `json:"specialization"`
	Status         string    `json:"status"`
}

type TimeStamp struct {
	StartTime string `json:"start_time" gorm:"not null"` // Use int64 to store timestamps
	EndTime   string `json:"end_time" gorm:"not null"`
}
