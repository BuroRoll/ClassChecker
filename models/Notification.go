package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Receiver uint
	Type     string
	Data     string
	IsRead   bool
}
