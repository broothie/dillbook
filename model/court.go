package model

import "gorm.io/gorm"

type Court struct {
	gorm.Model
	Name     string
	Bookings []Booking
}
