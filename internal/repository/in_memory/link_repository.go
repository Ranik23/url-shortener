package inmemory

import (
	"context"
	"log/slog"

	"github.com/Ranik23/url-shortener/internal/repository"
)



type linkRepository struct {
	links map[string]string
	logger *slog.Logger
}


func NewLinkRepostiory(logger *slog.Logger) repository.LinkRepository{
	return &linkRepository{
		links: make(map[string]string),
		logger: logger,
	}
}

func (l *linkRepository) CreateLink(ctx context.Context, default_link string, shortened_link string) error {
	l.logger.Info("CreateLIink")
	l.links[default_link] = shortened_link
	return nil
}

// DeleteLink implements repository.linkRepository.
func (l *linkRepository) DeleteLink(ctx context.Context, default_link string) error {
	delete(l.links, default_link)
	return nil
}

// GetDefaultLink implements repository.linkRepository.
func (l *linkRepository) GetDefaultLink(ctx context.Context, shortened_link string) (default_link string, err error) {

	for key, value := range l.links {
		if value == shortened_link {
			return key, nil
		}
	}
	return "", repository.ErrNotFound

}

// GetShortenedLink implements repository.linkRepository.
func (l *linkRepository) GetShortenedLink(ctx context.Context, default_link string) (shortened_link string, err error) {
	value, ok := l.links[default_link]
	if !ok {
		return "", repository.ErrNotFound
	}
	return value, nil
}