package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func New(log *zap.SugaredLogger) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return &pgx.Conn{}, err
	}

	return conn, err
}
