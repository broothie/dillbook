package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	gorm.Model
	ID string `gorm:"type:uuid;"`
}

func (b *Base) BeforeCreate(*gorm.DB) error {
	b.ID = uuid.NewString()
	return nil
}

type Location struct {
	Base
	Name        string `gorm:"not null;"`
	Address     string `gorm:"not null;"`
	Description string
	Courts      []Court
}

type Court struct {
	Base
	Name       string `gorm:"not null;"`
	LocationID string `gorm:"not null;"`
	Location   Location
	Bookings   []Booking
}

type Booking struct {
	Base
	Name        string
	Email       string
	PhoneNumber string
	Time        time.Time     `gorm:"not null;"`
	Duration    time.Duration `gorm:"not null;"`
	CourtID     string
	Court       Court
}
