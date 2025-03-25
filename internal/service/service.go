package service

import (
	"context"
	"errors"
)



var (
	ErrNotFound = errors.New("not found")
	ErrEmptyURL = errors.New("short url cannot be empty")	
)

type Service interface {
	LinkService
	StatService
	UserService
}

type LinkService interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	DeleteShortURL(ctx context.Context, shortURL string) 	error
	ResolveShortURL(ctx context.Context, shortURL string) (string, error) 
}

type StatService interface {
	GetStats(ctx context.Context, shortURL string) (any, error)
}

type UserService interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, username string) error
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

