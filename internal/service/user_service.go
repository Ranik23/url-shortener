package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Ranik23/url-shortener/internal/repository"
)


type UserService interface {
	CreateUser(ctx context.Context, username string) error
	DeleteUser(ctx context.Context, username string) error
}

type userService struct {
	txManager 	repository.TxManager
	userRepo 	repository.UserRepository
	logger		*slog.Logger
}

func (u *userService) CreateUser(ctx context.Context, username string) error {
	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := u.userRepo.CreateUser(txCtx, username); err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				return ErrAlreadyExists
			}
			return err
		}
		return nil
	})
}

func (u *userService) DeleteUser(ctx context.Context, username string) error {
	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := u.userRepo.DeleteUser(txCtx, username); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}
		return nil
	})
}

func NewUserService(userRepo repository.UserRepository, txManager repository.TxManager, logger *slog.Logger) UserService {
	return &userService{
		txManager: txManager,
		userRepo: userRepo,
		logger: logger,
	}
}
