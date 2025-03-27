package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Ranik23/url-shortener/internal/repository"
	servicehelpers "github.com/Ranik23/url-shortener/internal/libs/service_helpers"
)

type LinkService interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	DeleteShortURL(ctx context.Context, shortURL string) 	error
	ResolveShortURL(ctx context.Context, shortURL string) (string, error) 
}

type linkService struct {
	linkRepo 	repository.LinkRepository
	txManager 	repository.TxManager
}

func (l *linkService) CreateShortURL(ctx context.Context, default_link string) (string, error) {
	if strings.TrimSpace(default_link) == "" {
		return "", ErrEmptyURL
	}

	var short_link string

	err := l.txManager.Do(ctx, func(txCtx context.Context) error {

		shortened_link, err := l.linkRepo.GetShortenedLink(txCtx, default_link)
		if err != nil && !errors.Is(err, repository.ErrNotFound){
			return ErrInternal
		}

		short_link, err = servicehelpers.GenereateShortenedLink(default_link)
		if err != nil {
			return ErrInternal
		}

		if errors.Is(err, repository.ErrNotFound) {
			if err := l.linkRepo.CreateLink(txCtx, default_link, short_link); err != nil {
				return ErrInternal
			}
			return nil
		}

		short_link = shortened_link
		return nil
	})

	if err != nil {
		return "", ErrInternal
	}

	return short_link, nil
}

func (l *linkService) DeleteShortURL(ctx context.Context, shortURL string) error {
	if strings.TrimSpace(shortURL) == "" {
		return ErrEmptyURL
	}
	return l.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := l.linkRepo.DeleteLink(txCtx, shortURL); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return ErrInternal
		}
		return nil
	})
}


func (l *linkService) ResolveShortURL(ctx context.Context, shortURL string) (string, error) {
	if strings.TrimSpace(shortURL) == "" {
		return "", ErrEmptyURL
	}

	var default_link string
	
	err := l.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		default_link, err = l.linkRepo.GetDefaultLink(txCtx, shortURL)
		if err != nil {
			return ErrInternal
		}
		return nil
	})

	if err != nil {
		return "", ErrInternal
	}

	return default_link, nil
}

func NewLinkService(linkRepo repository.LinkRepository, txManager repository.TxManager) LinkService {
	return &linkService{
		linkRepo: linkRepo,
		txManager: txManager,
	}
}
