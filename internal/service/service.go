package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/p1xart/shortner-service/internal/repo/repoerrors"
	"go.uber.org/zap"
)

type Shortner interface {
	ReduceLink(ctx context.Context, srcLink, reduceLink string) error
	GetShortBySource(ctx context.Context, srcLink string) (string, error)
	GetSourceByShort(ctx context.Context, shortLink string) (string, error)
	IncrementVisitsByShort(ctx context.Context, shortLink string) error
}

type ShortnerService struct {
	log *zap.SugaredLogger

	repo Shortner
}

func NewService(log *zap.SugaredLogger, repo Shortner) *ShortnerService {
	log = log.With("component", "service")

	return &ShortnerService{
		log: log,

		repo: repo,
	}
}

func (s *ShortnerService) ReduceLink(ctx context.Context, srcLink string) (string, error) {
	var reduceLink string
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	reduceLink, err := s.repo.GetShortBySource(ctx, srcLink)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			s.log.Debugln("link not exists")
		} else {
			s.log.Errorln("error get short link by source link", zap.Error(err))
			return "", err
		}
	} else {
		s.log.Debugln("link already exists; returning", zap.String("source link", srcLink), zap.String("short link", reduceLink))
		return reduceLink, nil
	}

	repeats := 0
	for {
		reduceLink = randStringRunes(5, letterRunes)

		err := s.repo.ReduceLink(ctx, srcLink, reduceLink)
		if err != nil {
			if errors.Is(err, repoerrors.ErrAlreadyExists) {
				s.log.Debugln("short link already exists; regeneration...")
				repeats++
				if repeats < 2 {
					continue
				}
				return "", ErrLinkExists
			}
			s.log.Errorln("failed to create short link", zap.String("source link", srcLink), zap.Error(err))
			return "", err
		}
		break
	}

	s.log.Debugln("created new link", zap.String("source link", srcLink), zap.String("short link", reduceLink))
	return reduceLink, nil
}

func (s *ShortnerService) GetSourceByShort(ctx context.Context, shortLink string) (string, error) {
	sourceLink, err := s.repo.GetSourceByShort(ctx, shortLink)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return "", ErrLinkNotFound
		}
	}

	return sourceLink, nil
}

func (s *ShortnerService) IncrementVisitsByShort(ctx context.Context, shortLink string) error {
	if err := s.repo.IncrementVisitsByShort(ctx, shortLink); err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return ErrLinkNotFound
		}
		return err
	}

	return nil
}

func randStringRunes(lenght int, letterRunes []rune) string {
	b := make([]rune, lenght)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
