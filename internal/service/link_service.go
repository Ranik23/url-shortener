package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/Ranik23/url-shortener/internal/repository"
	"log/slog"
	"strings"
)

type LinkService interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	DeleteShortURL(ctx context.Context, shortURL string) error
	ResolveShortURL(ctx context.Context, shortURL string) (string, error)
}

type linkService struct {
	linkRepo  repository.LinkRepository
	txManager repository.TxManager
	logger    *slog.Logger
}

func (l *linkService) CreateShortURL(ctx context.Context, default_link string) (string, error) {
	if strings.TrimSpace(default_link) == "" {
		l.logger.Warn("CreateShortURL: empty URL provided")
		return "", ErrEmptyURL
	}

	var short_link string

	err := l.txManager.Do(ctx, func(txCtx context.Context) error {
		l.logger.Info("CreateShortURL: generating short link", slog.String("original_url", default_link))

		shortened_link, getErr := l.linkRepo.GetShortenedLink(txCtx, default_link)
		if getErr != nil && !errors.Is(getErr, repository.ErrNotFound) {
			l.logger.Error("CreateShortURL: error fetching shortened link", slog.String("error", getErr.Error()))
			return ErrInternal
		}

		if errors.Is(getErr, repository.ErrNotFound) {
			var genErr error
			short_link, genErr = l.generateShortenedLink(default_link)
			if genErr != nil {
				l.logger.Error("CreateShortURL: error generating short link", slog.String("error", genErr.Error()))
				return ErrInternal
			}
			l.logger.Info("CreateShortURL: link not found, creating new entry", slog.String("short_link", short_link))
			if err := l.linkRepo.CreateLink(txCtx, default_link, short_link); err != nil {
				l.logger.Error("CreateShortURL: error creating link", slog.String("error", err.Error()))
				return ErrInternal
			}
		} else {
			short_link = shortened_link
		}
		return nil
	})

	if err != nil {
		l.logger.Error("CreateShortURL: transaction failed", slog.String("error", err.Error()))
		return "", ErrInternal
	}

	l.logger.Info("CreateShortURL: successfully created short link", slog.String("short_link", short_link))
	return short_link, nil
}

func (l *linkService) DeleteShortURL(ctx context.Context, shortURL string) error {
	if strings.TrimSpace(shortURL) == "" {
		l.logger.Warn("DeleteShortURL: empty URL provided")
		return ErrEmptyURL
	}

	err := l.txManager.Do(ctx, func(txCtx context.Context) error {
		l.logger.Info("DeleteShortURL: attempting to delete", slog.String("short_url", shortURL))

		if err := l.linkRepo.DeleteLink(txCtx, shortURL); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				l.logger.Warn("DeleteShortURL: link not found", slog.String("short_url", shortURL))
				return ErrNotFound
			}
			l.logger.Error("DeleteShortURL: error deleting link", slog.String("error", err.Error()))
			return ErrInternal
		}

		l.logger.Info("DeleteShortURL: successfully deleted", slog.String("short_url", shortURL))
		return nil
	})

	if err != nil {
		l.logger.Error("DeleteShortURL: transaction failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (l *linkService) ResolveShortURL(ctx context.Context, shortURL string) (string, error) {
	if strings.TrimSpace(shortURL) == "" {
		l.logger.Warn("ResolveShortURL: empty URL provided")
		return "", ErrEmptyURL
	}

	var default_link string

	err := l.txManager.Do(ctx, func(txCtx context.Context) error {
		l.logger.Info("ResolveShortURL: resolving", slog.String("short_url", shortURL))

		var err error
		default_link, err = l.linkRepo.GetDefaultLink(txCtx, shortURL)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return repository.ErrNotFound
			}
			l.logger.Error("ResolveShortURL: error resolving link", slog.String("error", err.Error()))
			return ErrInternal
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", repository.ErrNotFound
		}
		l.logger.Error("ResolveShortURL: transaction failed", slog.String("error", err.Error()))
		return "", ErrInternal
	}

	l.logger.Info("ResolveShortURL: successfully resolved", slog.String("original_url", default_link))
	return default_link, nil
}

func NewLinkService(linkRepo repository.LinkRepository, txManager repository.TxManager, logger *slog.Logger) LinkService {
	return &linkService{
		linkRepo:  linkRepo,
		txManager: txManager,
		logger:    logger,
	}
}

func (l *linkService) generateShortenedLink(defaultLink string) (string, error) {
	if defaultLink == "" {
		return "", errors.New("empty URL not allowed")
	}
	hash := sha256.Sum256([]byte(defaultLink))
	shortURL := base64.RawURLEncoding.EncodeToString(hash[:])[:8]
	return shortURL, nil
}
