package models

import "gorm.io/gorm"

type BookingTime struct {
	gorm.Model
	BookingClassID uint
	Time           string
	IsEnd          bool
	IsSuccess      bool
}
