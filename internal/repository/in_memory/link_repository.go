package inmemory

import (
	"context"
	"github.com/Ranik23/url-shortener/internal/repository"
)

type inMemoryRepository struct {
	links map[string]string
	users map[string]int
}

func (u *inMemoryRepository) CreateUser(ctx context.Context, username string) error {
	u.users[username] += 1
	return nil
}

// DeleteUser implements repository.UserRepository.
func (u *inMemoryRepository) DeleteUser(ctx context.Context, username string) error {
	delete(u.users, username)
	return nil
}

// UserExists implements repository.UserRepository.
func (u *inMemoryRepository) UserExists(ctx context.Context, username string) (exists bool, err error) {
	return u.users[username] > 0, nil
}


// CreateLink implements repository.inMemoryRepository.
func (l *inMemoryRepository) CreateLink(ctx context.Context, default_link string, shortened_link string) error {
	l.links[default_link] = shortened_link
	return nil
}

// DeleteLink implements repository.inMemoryRepository.
func (l *inMemoryRepository) DeleteLink(ctx context.Context, default_link string) error {
	delete(l.links, default_link)
	return nil
}

// GetDefaultLink implements repository.inMemoryRepository.
func (l *inMemoryRepository) GetDefaultLink(ctx context.Context, shortened_link string) (default_link string, err error) {

	for key, value := range l.links {
		if value == shortened_link {
			return key, nil
		}
	}
	return "", repository.ErrNotFound

}

// GetShortenedLink implements repository.inMemoryRepository.
func (l *inMemoryRepository) GetShortenedLink(ctx context.Context, default_link string) (shortened_link string, err error) {
	value, ok := l.links[default_link]
	if !ok {
		return "", repository.ErrNotFound
	}
	return value, nil
}

func NewInMemoryRepository() repository.Repository {
	return &inMemoryRepository{}
}
