package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/Ranik23/url-shortener/internal/repository"
	repoMock "github.com/Ranik23/url-shortener/internal/repository/mock"
	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

// computeExpectedShortLink рассчитывает значение короткой ссылки, тот же generateShortenedLink
func computeExpectedShortLink(originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("empty URL not allowed")
	}
	hash := sha256.Sum256([]byte(originalURL))
	shortURL := base64.RawURLEncoding.EncodeToString(hash[:])[:8]
	return shortURL, nil
}

func TestCreateShortURL_Empty(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	short, err := svc.CreateShortURL(context.Background(), "   ")
	assert.Equal(t, "", short)
	assert.ErrorIs(t, err, ErrEmptyURL)
}

func TestCreateShortURL_Existing(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	originalURL := "http://example.com"
	existingShort := "exist123"

	// При попытке получить сокращённый URL репозиторий возвращает существующее значение
	linkRepo.
		On("GetShortenedLink", mock.Anything, originalURL).
		Return(existingShort, nil)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(nil).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			_ = fn(context.Background())
		})

	short, err := svc.CreateShortURL(context.Background(), originalURL)
	assert.NoError(t, err)
	assert.Equal(t, existingShort, short)

	linkRepo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}

func TestCreateShortURL_NewEntry(t *testing.T) {
	// Тестируем создание новой записи
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	originalURL := "http://example.com"

	expectedShort, err := computeExpectedShortLink(originalURL)
	require.NoError(t, err)

	// Ситуация, когда для данного URL не найдена запись
	linkRepo.
		On("GetShortenedLink", mock.Anything, originalURL).
		Return("", repository.ErrNotFound)

	linkRepo.
		On("CreateLink", mock.Anything, originalURL, expectedShort).
		Return(nil)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(nil).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			_ = fn(context.Background())
		})

	short, err := svc.CreateShortURL(context.Background(), originalURL)
	assert.NoError(t, err)
	assert.Equal(t, expectedShort, short)

	linkRepo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}

func TestDeleteShortURL_Empty(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	err := svc.DeleteShortURL(context.Background(), "   ")
	assert.ErrorIs(t, err, ErrEmptyURL)
}

func TestDeleteShortURL_NotFound(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	shortURL := "short123"

	// Ситуация, когда удаление возвращает ErrNotFound
	linkRepo.
		On("DeleteLink", mock.Anything, shortURL).
		Return(repository.ErrNotFound)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	err := svc.DeleteShortURL(context.Background(), shortURL)
	assert.ErrorIs(t, err, ErrNotFound)
	linkRepo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}

func TestDeleteShortURL_Success(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	shortURL := "short123"

	linkRepo.
		On("DeleteLink", mock.Anything, shortURL).
		Return(nil)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(nil).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			_ = fn(context.Background())
		})

	err := svc.DeleteShortURL(context.Background(), shortURL)
	assert.NoError(t, err)
	linkRepo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}

func TestResolveShortURL_Empty(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	orig, err := svc.ResolveShortURL(context.Background(), "   ")
	assert.Equal(t, "", orig)
	assert.ErrorIs(t, err, ErrEmptyURL)
}

func TestResolveShortURL_Error(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	shortURL := "short123"
	// Ошибка при получении оригинального URL
	linkRepo.
		On("GetDefaultLink", mock.Anything, shortURL).
		Return("", errors.New("db error"))

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	orig, err := svc.ResolveShortURL(context.Background(), shortURL)
	assert.Equal(t, "", orig)
	assert.ErrorIs(t, err, ErrInternal)
	linkRepo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}

func TestResolveShortURL_Success(t *testing.T) {
	txManager := repoMock.NewTxManager(t)
	linkRepo := repoMock.NewLinkRepository(t)
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	svc := NewLinkService(linkRepo, txManager, logger)

	shortURL := "short123"
	originalURL := "http://example.com"

	linkRepo.
		On("GetDefaultLink", mock.Anything, shortURL).
		Return(originalURL, nil)

	txManager.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Return(nil).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			_ = fn(context.Background())
		})

	orig, err := svc.ResolveShortURL(context.Background(), shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, orig)
	linkRepo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}
