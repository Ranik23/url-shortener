package service

import (
	"context"

	"github.com/Ranik23/url-shortener/internal/repository"
)

type userService struct {
	txManager repository.TxManager
	userRepo repository.UserRepository
}

func (u *userService) CreateUser(ctx context.Context, username string) error {
	panic("unimplemented")
}

func (u *userService) DeleteUser(ctx context.Context, username string) error {
	panic("unimplemented")
}

func NewUserService(userRepo repository.UserRepository, txManager repository.TxManager) UserService {
	return &userService{
		txManager: txManager,
		userRepo: userRepo,
	}
}
