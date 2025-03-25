package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/Ranik23/url-shortener/internal/service/utils"
	"github.com/jackc/pgx/v5"
)

type linkService struct {
	linkRepo 	repository.LinkRepository
	txManager 	repository.TxManager
}

func (l *linkService) CreateShortURL(ctx context.Context, default_link string) (string, error) {

	var short_link string

	err := l.txManager.WithTx(ctx, pgx.Serializable, pgx.ReadWrite, func(txCtx context.Context) error {

		shortened_link, err := l.linkRepo.GetShortenedLink(txCtx, default_link)
		if err != nil && !errors.Is(err, repository.ErrNotFound){
			return err
		}

		short_link, err = utils.GenereateShortenedLink(default_link)
		if err != nil {
			return err
		}

		if errors.Is(err, repository.ErrNotFound) {
			if err := l.linkRepo.CreateLink(txCtx, default_link, short_link); err != nil {
				return err
			}
			return nil
		}

		short_link = shortened_link
		return nil
	})

	if err != nil {
		return "", err
	}

	return short_link, nil
}

func (l *linkService) DeleteShortURL(ctx context.Context, shortURL string) error {
	if strings.TrimSpace(shortURL) == "" {
		return ErrEmptyURL
	}
	return l.txManager.WithTx(ctx, pgx.Serializable, pgx.ReadWrite, func(txCtx context.Context) error {
		if err := l.linkRepo.DeleteLink(txCtx, shortURL); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}
		return nil
	})
}


func (l *linkService) ResolveShortURL(ctx context.Context, shortURL string) (string, error) {
	if strings.TrimSpace(shortURL) == "" {
		return "", ErrEmptyURL
	}

	var default_link string
	
	err := l.txManager.WithTx(ctx, pgx.Serializable, pgx.ReadOnly, func(txCtx context.Context) error {
		var err error
		default_link, err = l.linkRepo.GetDefaultLink(txCtx, shortURL)
		return err
	})

	if err != nil {
		return "", fmt.Errorf("failed to resolve short URL: %w", err)
	}

	return default_link, nil
}

func NewLinkService(linkRepo repository.LinkRepository, txManager repository.TxManager) LinkService {
	return &linkService{
		linkRepo: linkRepo,
		txManager: txManager,
	}
}
