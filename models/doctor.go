package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	Name           string `json:"name" gorm:"not null"`
	Email          string `json:"email" gorm:"unique;not null"`
	Specialization string `json:"specialization" gorm:"not null"`
}
