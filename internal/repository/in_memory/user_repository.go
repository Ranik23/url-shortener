package inmemory

import (
	"context"
	"log/slog"

	"github.com/Ranik23/url-shortener/internal/repository"
)


type userRepository struct {
	users map[string]int
	logger *slog.Logger 
}

func NewUserRepository(logger *slog.Logger) repository.UserRepository {
	return &userRepository{
		logger: logger,
		users: make(map[string]int),
	}
}


func (u *userRepository) CreateUser(ctx context.Context, username string) error {
	u.logger.Info("CreateUser")
	u.users[username] += 1
	return nil

}

// DeleteUser implements repository.UserRepository.
func (u *userRepository) DeleteUser(ctx context.Context, username string) error {
	delete(u.users, username)
	return nil
}


//serExists implements repository.UserRepository.
func (u *userRepository) UserExists(ctx context.Context, username string) (exists bool, err error) {
	return u.users[username] > 0, nil
}
