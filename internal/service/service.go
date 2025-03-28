package service

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrEmptyURL = errors.New("short url cannot be empty")	
	ErrInternal = errors.New("internal error")
)

type Service interface {
	LinkService
	StatService
	UserService
}


type service struct {
	LinkService
	StatService
	UserService
}


func NewService(linkService LinkService, statService StatService, userService UserService) Service {
	return &service{
		LinkService: linkService,
		StatService: statService,
		UserService: userService,
	}
}

