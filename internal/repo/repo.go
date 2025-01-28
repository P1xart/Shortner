package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/p1xart/shortner-service/internal/repo/codes"
	"github.com/p1xart/shortner-service/internal/repo/repoerrors"
	"go.uber.org/zap"
)

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repo struct {
	log *zap.SugaredLogger

	pg Pool
}

func NewRepo(log *zap.SugaredLogger, postgres Pool) *Repo {
	log = log.With("component", "repo")

	return &Repo{
		log: log,

		pg: postgres,
	}
}

func (r *Repo) ReduceLink(ctx context.Context, srcLink, reduceLink string) error {
	query := "INSERT INTO links(src_link, short_link) VALUES($1, $2);"

	_, err := r.pg.Exec(ctx, query, srcLink, reduceLink)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == codes.UniqueConstraintCode {
				r.log.Debugln("link already exists", zap.String("short link", reduceLink), zap.Error(err))
				return repoerrors.ErrAlreadyExists
			}
		}
		r.log.Errorln("error to reduce link", zap.String("short link", reduceLink), zap.Error(err))
		return err
	}

	return err
}

func (r *Repo) GetShortBySource(ctx context.Context, srcLink string) (string, error) {
	var shortLink string
	query := "SELECT short_link FROM links WHERE src_link=$1;"

	if err := r.pg.QueryRow(ctx, query, srcLink).Scan(&shortLink); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debugln("link not exists", zap.String("source link", srcLink), zap.Error(err))
			return "", repoerrors.ErrNotFound
		}
		r.log.Errorln("failed get short link by source link", zap.String("source link", srcLink), zap.Error(err))
		return "", err
	}

	return shortLink, nil
}

func (r *Repo) GetSourceByShort(ctx context.Context, shortLink string) (string, error) {
	var srcLink string
	query := "SELECT src_link FROM links WHERE short_link=$1;"

	if err := r.pg.QueryRow(ctx, query, shortLink).Scan(&srcLink); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debugln("link not exists", zap.String("short link", shortLink), zap.Error(err))
			return "", repoerrors.ErrNotFound
		}
		r.log.Errorln("failed get source link by short link", zap.String("short link", srcLink), zap.String("source link", srcLink), zap.Error(err))
		return "", err
	}

	return srcLink, nil
}

func (r *Repo) IncrementVisitsByShort(ctx context.Context, shortLink string) error {
	query := "UPDATE links SET visits=visits+1 WHERE short_link=$1"

	tag, err := r.pg.Exec(ctx, query, shortLink)
	if err != nil {
		r.log.Errorln("failed to increment visit", zap.String("short link", shortLink), zap.Error(err))
		return err
	}
	if tag.RowsAffected() == 0 {
		r.log.Debugln("link not exists", zap.String("short link", shortLink), zap.Error(err))
		return repoerrors.ErrNotFound
	}

	return nil
}
