package repository

import (
	"context"
	"errors"
)


var (
	ErrNotFound = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)


type LinkRepository interface {
	CreateLink(ctx context.Context, default_link string, shortened_link string) 	error
	GetDefaultLink(ctx context.Context, shortened_link string) 						(default_link string, err error)
	GetShortenedLink(ctx context.Context, default_link string) 						(shortened_link string, err error)

	DeleteLink(ctx context.Context, default_link string)							error
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, username string) error 
	UserExists(ctx context.Context, username string) (exists bool, err error)
}

type Repository interface {
	LinkRepository
	UserRepository
}


type repository struct {
	UserRepository
	LinkRepository
}

func NewRepository(userRepo UserRepository, linkRepo LinkRepository) Repository {
	return &repository{
		UserRepository: userRepo,
		LinkRepository: linkRepo,
	}
}
