package model

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	ID          string `gorm:"type:uuid;default:uuid_generate_v4();"`
	Name        string `gorm:"not null;"`
	Address     string `gorm:"not null;"`
	Description string
	Courts      []Court
}

type Court struct {
	gorm.Model
	ID         string `gorm:"type:uuid;default:uuid_generate_v4();"`
	Name       string `gorm:"not null;"`
	LocationID string `gorm:"not null;"`
	Location   Location
	Bookings   []Booking
}

type Booking struct {
	gorm.Model
	ID          string `gorm:"type:uuid;default:uuid_generate_v4();"`
	Name        string
	Email       string
	PhoneNumber string
	Time        time.Time     `gorm:"not null;"`
	Duration    time.Duration `gorm:"not null;"`
	CourtID     string
	Court       Court
}
