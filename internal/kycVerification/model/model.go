package model

import (
	"time"

	"gorm.io/gorm"
)

type PANVerification struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        string `gorm:"index"`
	PAN           string
	Consent       bool
	RequestReason string
	Status        bool `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type RPDVerification struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	TraceID   string `gorm:"index;unique"`
	Status    string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
