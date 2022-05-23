package repository

import (
	"ClassChecker/models"
	"gorm.io/gorm"
)

type Classes interface {
	GetClassesTimes() []models.BookingTime
	SaveClassesTimes(models.BookingTime) error
	GetClassesWithTime() []UserBooking
	SaveClass(times UserBooking) error
	SaveNotification(notification models.ClassNotification) models.ClassNotification
}

type Repository struct {
	Classes
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Classes: NewClassesPostgres(db),
	}
}
