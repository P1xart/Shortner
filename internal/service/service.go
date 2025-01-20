package service

import (
	"context"

	"go.uber.org/zap"
)

type Shortner interface {
	ReduceLink(ctx context.Context, link string) (string, error)
}

type ShortnerService struct {
	log *zap.SugaredLogger

	repo Shortner
}

func NewService(log *zap.SugaredLogger, repo Shortner) *ShortnerService {
	log = log.With("component", "repo")

	return &ShortnerService{
		log: log,

		repo: repo,
	}
}

func (s *ShortnerService) ReduceLink(ctx context.Context, link string) (string, error) {
	reduceLink, err := s.repo.ReduceLink(ctx, link)
	if err != nil {
		s.log.Error("error create reduce link", zap.Error(err))
		return "", err
	}

	return reduceLink, err
}
