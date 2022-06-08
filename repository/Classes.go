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
	c.db.
		Where("to_timestamp(CONCAT( substr(time, 1, length(time) - 2), ' ', (substr(time, length(time), length(time))::INTEGER * time '03:00')), 'YYYY/MM/DD HH24:MI:SS') < NOW() AND is_end=false").
		Find(&bookings)
	return bookings
}

func (c ClassesPostgres) SaveClassesTimes(times models.BookingTime) error {
	result := c.db.Model(&times).Where("id = ?", times.ID).Update("is_end", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type UserBooking struct {
	ID               uint
	UserID           uint
	ClassID          uint
	Status           string
	MentiId          uint
	ClassDataName    string `json:"class_data_name"`
	MentiFirstName   string `json:"menti_first_name"`
	MentiSecondName  string `json:"menti_second_name"`
	MentorFirstName  string `json:"mentor_first_name"`
	MentorSecondName string `json:"mentor_second_name"`
	LessonCount      uint   `json:"lesson_count"`

	Time []models.BookingTime `gorm:"foreignKey:BookingClassID;references:ID"`
}

func (c ClassesPostgres) GetClassesWithTime() []UserBooking {
	var bookings []UserBooking
	c.db.
		Table("user_classes").
		Select("*").
		Preload("Time").
		Joins("LEFT JOIN (SELECT id AS class_data_id, class_name AS class_data_name FROM classes) AS class_data ON class_data_id = class_id").
		Joins("LEFT JOIN (SELECT id as menti_data_id, first_name AS menti_first_name, second_name AS menti_second_name FROM users) AS menti_data ON menti_data_id = menti_id").
		Joins("LEFT JOIN (SELECT id AS mentor_data_id, first_name AS mentor_first_name, second_name AS mentor_second_name FROM users) AS mentor_data ON mentor_data_id = user_classes.user_id").
		Joins("LEFT JOIN (SELECT COUNT(id) AS lesson_count, booking_class_id AS booking_class_data_id FROM booking_times GROUP BY booking_class_id) as lcbcdi ON booking_class_data_id = id").
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

func (c ClassesPostgres) GetClassData(classId uint) UserBooking {
	var classDat UserBooking
	c.db.Table("user_classes").
		Select("*").
		Joins("LEFT JOIN (SELECT id AS class_data_id, class_name AS class_data_name FROM classes) AS class_data ON class_data_id = class_id").
		Joins("LEFT JOIN (SELECT id as menti_data_id, first_name AS menti_first_name, second_name AS menti_second_name FROM users) AS menti_data ON menti_data_id = menti_id").
		Joins("LEFT JOIN (SELECT id AS mentor_data_id, first_name AS mentor_first_name, second_name AS mentor_second_name FROM users) AS mentor_data ON mentor_data_id = user_classes.user_id").
		Where("id = ?", classId).
		Find(&classDat)
	return classDat
}
