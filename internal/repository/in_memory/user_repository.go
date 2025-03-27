package inmemory

import (
	"context"
	"github.com/Ranik23/url-shortener/internal/repository"
)

type userRepository struct {
	users map[string]int
}

// CreateUser implements repository.UserRepository.
func (u *userRepository) CreateUser(ctx context.Context, username string) error {
	u.users[username] += 1
	return nil
}

// DeleteUser implements repository.UserRepository.
func (u *userRepository) DeleteUser(ctx context.Context, username string) error {
	delete(u.users, username)
	return nil
}

// UserExists implements repository.UserRepository.
func (u *userRepository) UserExists(ctx context.Context, username string) (exists bool, err error) {
	return u.users[username] > 0, nil
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}
