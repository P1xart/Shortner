package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type PgxPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Repo struct {
	log *zap.SugaredLogger

	postgres PgxPool
}

func NewRepo(log *zap.SugaredLogger, postgres PgxPool) *Repo {
	return &Repo{
		log: log,

		postgres: postgres,
	}
}

func (r *Repo) ReduceLink(ctx context.Context, link string) (string, error) {
	// reduceLink, err := r.postgres.Query(ctx, link)
	return "New link", nil
}
