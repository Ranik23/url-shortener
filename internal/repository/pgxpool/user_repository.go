package pgxpool

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	ctxManager 	repository.CtxManager
	settings 	repository.Settings
}

func (p *userRepository) CreateUser(ctx context.Context, username string) error {
	tr := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Transaction().(pgx.Tx)


	sql, args, err := sq.Insert("users").
		Columns("username").
		Values(username).
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

func (p *userRepository) DeleteUser(ctx context.Context, username string) error {
	tr 	 := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	sql, args, err := sq.Delete("users").
		Where(sq.Eq{"username": username}).
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

func (p *userRepository) UserExists(ctx context.Context, username string) (bool, error) {
	tr 	 := p.ctxManager.ByKey(ctx, p.settings.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	sql, args, err := sq.Select("1").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()

	if err != nil {
		return false, err
	}

	var exists bool
	if err = exec.QueryRow(ctx, sql, args).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

func NewUserRepository(ctxManager repository.CtxManager, settings repository.Settings) repository.UserRepository {
	return &userRepository{
		ctxManager: ctxManager,
		settings: settings,
	}
}
