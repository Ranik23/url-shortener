//go:build integration

package integration

import (
	"context"
	"errors"
	"github.com/Ranik23/url-shortener/internal/repository"
)

func (s *TestSuite) TestCreateAndResolveShortURL() {
	ctx := context.Background()
	originalURL := "http://example.com"

	// Создаем короткую ссылку
	short, err := s.svc.CreateShortURL(ctx, originalURL)
	s.Require().NoError(err)
	s.Require().NotEmpty(short)

	// Делаю короткую ссылку обратно в оригинальный URL
	resolved, err := s.svc.ResolveShortURL(ctx, short)
	s.Require().NoError(err)
	s.Require().Equal(originalURL, resolved)
}

func (s *TestSuite) TestDeleteShortURL() {
	ctx := context.Background()
	originalURL := "http://example.com"

	short, err := s.svc.CreateShortURL(ctx, originalURL)
	s.Require().NoError(err)
	s.Require().NotEmpty(short)

	err = s.svc.DeleteShortURL(ctx, short)
	s.Require().NoError(err)

	_, err = s.svc.ResolveShortURL(ctx, short)
	s.Require().Error(err)
}

type failingLinkRepo struct {
	real repository.LinkRepository
}

func (f *failingLinkRepo) CreateLink(ctx context.Context, defaultLink, shortenedLink string) error {
	return f.real.CreateLink(ctx, defaultLink, shortenedLink)
}

func (f *failingLinkRepo) GetDefaultLink(ctx context.Context, shortenedLink string) (string, error) {
	return f.real.GetDefaultLink(ctx, shortenedLink)
}

func (f *failingLinkRepo) GetShortenedLink(ctx context.Context, defaultLink string) (string, error) {
	return f.real.GetShortenedLink(ctx, defaultLink)
}

func (f *failingLinkRepo) DeleteLink(ctx context.Context, shortURL string) error {
	// Симулируем ошибку удаления, чтобы вызвать rollback транзакции.
	return errors.New("simulated delete error")
}

func (s *TestSuite) TestRepositoryDeleteRollback() {
	ctx := context.Background()
	originalURL := "http://rollback-repo.com"

	//Создаем запись через сервис
	short, err := s.svc.CreateShortURL(ctx, originalURL)
	s.Require().NoError(err)
	s.Require().NotEmpty(short)

	// Проверяем, что запись присутствует
	resolved, err := s.svc.ResolveShortURL(ctx, short)
	s.Require().NoError(err)
	s.Require().Equal(originalURL, resolved)

	// Выполняем транзакцию вручную, вызывая метод репозитория DeleteLink напрямую
	// При этом после удаления искусственно возвращаем ошибку, чтобы произошел rollback
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		repo := s.linkRepo
		if err := repo.DeleteLink(txCtx, short); err != nil {
			return err
		}
		return errors.New("simulated error for rollback")
	})
	s.Require().Error(err, "An error is expected for initiating a rollback")

	// Проверяем, что запись осталась
	resolvedAfter, err := s.svc.ResolveShortURL(ctx, short)
	s.Require().NoError(err, "The record should be saved after the transaction is rolled back.")
	s.Require().Equal(originalURL, resolvedAfter)
}

func (s *TestSuite) TestTransactionCommit() {
	ctx := context.Background()
	originalURL := "http://commit-example.com"

	short, err := s.svc.CreateShortURL(ctx, originalURL)
	s.Require().NoError(err)
	s.Require().NotEmpty(short)

	// Проверяем, что запись существует
	resolved, err := s.svc.ResolveShortURL(ctx, short)
	s.Require().NoError(err)
	s.Require().Equal(originalURL, resolved)

	// Вызываем DeleteShortURL напрямую
	err = s.svc.DeleteShortURL(ctx, short)
	s.Require().NoError(err, "The transaction should be completed successfully, fixing the deletion.")

	// После коммита транзакции, попытка разрешить удаленную ссылку должна вернуть ошибку ErrNotFound
	_, err = s.svc.ResolveShortURL(ctx, short)
	s.Require().Error(err, "There should be no record after the commit.")
	s.Require().Equal(repository.ErrNotFound, err)
}
