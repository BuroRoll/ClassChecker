package repository

import (
	"ClassChecker/models"
	"gorm.io/gorm"
)

type ClassesPostgres struct {
	db *gorm.DB
}

func NewClassesPostgres(db *gorm.DB) *ClassesPostgres {
	return &ClassesPostgres{db: db}
}

func (c ClassesPostgres) GetClassesTimes() []models.BookingTime {
	var bookings []models.BookingTime
	c.db.Where("to_timestamp(CONCAT( substr(time, 1, length(time) - 2), ' ', (substr(time, length(time), length(time))::INTEGER * time '03:00')), 'YYYY/MM/DD HH24:MI:SS') < NOW()").
		Find(&bookings)
	return bookings
}

func (c ClassesPostgres) SaveClassesTimes(times models.BookingTime) error {
	result := c.db.Save(times)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type UserBooking struct {
	ID            uint
	UserID        uint
	ClassID       uint
	Status        string
	MentiId       uint
	ClassDataName string `json:"class_data_name"`

	Time []models.BookingTime `gorm:"foreignKey:BookingClassID;references:ID"`
}

func (c ClassesPostgres) GetClassesWithTime() []UserBooking {
	var bookings []UserBooking
	c.db.
		Table("user_classes").
		Select("*").
		Preload("Time").
		Joins("LEFT JOIN (SELECT id AS class_data_id, class_name AS class_data_name FROM classes) AS class_data ON class_data_id = class_id").
		Where("status = ?", "planned").
		Find(&bookings)
	return bookings
}

func (c ClassesPostgres) SaveClass(times UserBooking) error {
	result := c.db.Table("user_classes").
		Where("id = ?", times.ID).
		Update("status", times.Status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c ClassesPostgres) SaveNotification(notification models.ClassNotification) models.ClassNotification {
	c.db.Create(&notification)
	return notification
}
