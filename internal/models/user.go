package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string    `json:"name"`
	Email            string    `gorm:"uniqueIndex" json:"email"`
	Password         string    `json:"-"` 
	IsVerified       bool      `json:"is_verified" gorm:"default:false"`
	VerificationCode string    `json:"-"`
	OTPExpiresAt     time.Time `json:"-"`
}
