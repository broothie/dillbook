package model

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	CourtID uint
	Users   []User `gorm:"many2many:user_bookings"`
}
