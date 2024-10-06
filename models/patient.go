package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	Email string `json:"email" gorm:"unique;not null"`
	Name  string `json:"name" gorm:"not null"`
}
