package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Bookings []Booking `gorm:"many2many:user_bookings"`
}