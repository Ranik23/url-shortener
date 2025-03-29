package pgxpool

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

type linkRepository struct {
	ctxManager repository.CtxManager
	settings   repository.Settings
}

func (p *linkRepository) CreateLink(ctx context.Context, default_link string, shortened_link string) error {
	tr := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Transaction().(pgx.Tx)

	sql, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).Insert("links").
		Columns("default_link", "shortened_link").
		Values(default_link, shortened_link).
		ToSql()

	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

// GetDefaultLink implements repository.LinkRepository.
func (p *linkRepository) GetDefaultLink(ctx context.Context, shortened_link string) (default_link string, err error) {
	tr := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	sql, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).Select("default_link").
		From("links").
		Where(sq.Eq{"shortened_link": shortened_link}).
		ToSql()

	if err != nil {
		return "", err
	}

	if err = exec.QueryRow(ctx, sql, args...).Scan(&default_link); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", err
	}

	return default_link, nil
}

// GetShortenedLink implements repository.LinkRepository.
func (p *linkRepository) GetShortenedLink(ctx context.Context, default_link string) (shortened_link string, err error) {
	tr := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	sql, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).Select("shortened_link").
		From("links").
		Where(sq.Eq{"default_link": default_link}).
		ToSql()

	if err != nil {
		return "", err
	}

	if err = exec.QueryRow(ctx, sql, args...).Scan(&shortened_link); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", err
	}

	return shortened_link, nil
}

func (p *linkRepository) DeleteLink(ctx context.Context, shortLink string) error {
	tr := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	sql, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).Delete("links").
		Where(sq.Eq{"shortened_link": shortLink}).
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
func NewLinkRepository(ctxManager repository.CtxManager, settings repository.Settings) repository.LinkRepository {
	return &linkRepository{
		settings:   settings,
		ctxManager: ctxManager,
	}
}
