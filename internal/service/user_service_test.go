package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Ranik23/url-shortener/internal/repository"
	repoMock "github.com/Ranik23/url-shortener/internal/repository/mock"
	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)




func TestCreateUser_Success(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	userRepo := repoMock.NewUserRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewUserService(userRepo, txManager, logger)

	userRepo.
		On("CreateUser", mock.Anything, "anton").
		Return(nil)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(nil).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			_ = fn(context.Background())
		})

	err := svc.CreateUser(context.Background(), "anton")
	assert.Nil(t, err)
	
	txManager.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}


func TestCreateUser_AlreadyExists(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	userRepo := repoMock.NewUserRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewUserService(userRepo, txManager, logger)

	userRepo.
		On("CreateUser", mock.Anything, "anton").
		Return(repository.ErrAlreadyExists)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	err := svc.CreateUser(context.Background(), "anton")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrAlreadyExists)
}