package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Bookings []Booking `gorm:"many2many:user_bookings"`
}

type Location struct {
	gorm.Model
	Name   string
	Courts []Court
}

type Court struct {
	gorm.Model
	Name     string
	Bookings []Booking
}

type Booking struct {
	gorm.Model
	CourtID uint
	Users   []User `gorm:"many2many:user_bookings"`
}
