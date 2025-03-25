package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

type postgresLinkRepository struct {
	txManager repository.TxManager
}

// CreateLink implements repository.LinkRepository.
func (p *postgresLinkRepository) CreateLink(ctx context.Context, default_link string, shortened_link string) error {
	exec := p.txManager.GetExecutor(ctx)

	sql, args, err := sq.Insert("links").
		Columns("default_link", "shortened_link").
		Values(default_link, shortened_link).
		ToSql()

	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}

// GetDefaultLink implements repository.LinkRepository.
func (p *postgresLinkRepository) GetDefaultLink(ctx context.Context, shortened_link string) (default_link string, err error) {
	exec := p.txManager.GetExecutor(ctx)

	sql, args, err := sq.Select("default_link").
		From("links").
		Where(sq.Eq{"shortened_link": shortened_link}).
		ToSql()

	if err != nil {
		return "", err
	}

	if err = exec.QueryRow(ctx, sql, args).Scan(&default_link); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", err
	}

	return default_link, nil
}

// GetShortenedLink implements repository.LinkRepository.
func (p *postgresLinkRepository) GetShortenedLink(ctx context.Context, default_link string) (shortened_link string, err error) {
	exec := p.txManager.GetExecutor(ctx)

	sql, args, err := sq.Select("shortened_link").
		From("links").
		Where(sq.Eq{"default_link": default_link}).
		ToSql()

	if err != nil {
		return "", err
	}

	if err = exec.QueryRow(ctx, sql, args).Scan(&shortened_link); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", err
	}

	return shortened_link, nil
}

func (p *postgresLinkRepository) DeleteLink(ctx context.Context, default_link string) error {

	exec := p.txManager.GetExecutor(ctx)

	sql, args, err := sq.Delete("links").
		Where(sq.Eq{"default_link": default_link}).
		ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := exec.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func NewPostgresLinkRepository(txManager repository.TxManager) repository.LinkRepository {
	return &postgresLinkRepository{
		txManager: txManager,
	}
}
