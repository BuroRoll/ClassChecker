package models

import "gorm.io/gorm"

type ClassNotification struct {
	gorm.Model
	Receiver uint
	Type     string
	Data     string
	IsRead   bool
}
