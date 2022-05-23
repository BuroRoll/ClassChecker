package service

import (
	"ClassChecker/repository"
)

type Class interface {
	CheckClasses()
	CheckClassEnd()
}

type Service struct {
	Class
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Class: NewClassesService(repos.Classes),
	}
}
