package service

import (
	"context"
	"errors"
	"time"

	"github.com/ntduc6403/bookmark-management/internal/repository"
	"github.com/ntduc6403/bookmark-management/pkg/stringutils"
	"github.com/redis/go-redis/v9"
)

const (
	urlCodeLength = 7
	defaultURLExp = 24 * time.Hour
)

//go:generate mockery --name ShortenUrl -- filename shortenurl.go
type ShortenUrl interface {
	ShortenUrl(ctx context.Context, url string, expTime int) (string, error)
}

type shortenUrl struct {
	repo    repository.URLStorage
	genCode func(length int) (string, error)
}

func NewShortenUrl(repo repository.URLStorage) ShortenUrl {

	return &shortenUrl{
		repo:    repo,
		genCode: stringutils.GenerateCode,
	}
}

func (s *shortenUrl) ShortenUrl(ctx context.Context, url string, expTime int) (string, error) {

	urlCode, err := s.genCode(urlCodeLength)

	if err != nil {
		return " ", err
	}
	res, err := s.repo.GetURL(ctx, urlCode)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	if res != "" {
		return s.ShortenUrl(ctx, url, expTime)
	}

	expiration := defaultURLExp

	if expTime > 0 {
		expiration = time.Duration(expTime) * time.Second
	}

	err = s.repo.StoreURL(ctx, urlCode, url, expiration)
	if err != nil {
		return " ", err
	}
	return urlCode, nil
}
