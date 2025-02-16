package service

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"time"
	"url-shortener/internal/entity"
	"url-shortener/internal/lib/reflecthelper"
	"url-shortener/internal/repository/reperrors"
	"url-shortener/internal/service/serverrors"
)

type Hashing interface {
	EncodeInt64(numbers []int64) (string, error)
	DecodeInt64WithError(hash string) ([]int64, error)
}

type LinkRepository interface {
	CreateLink(ctx context.Context, sourceUrl string, expiresAt *time.Time) (int64, error)
	GetLinkBySourceUrl(ctx context.Context, sourceUrl string) (*entity.Link, error)
	GetLinkById(ctx context.Context, id int64) (*entity.Link, error)
	UpdateLastRequested(ctx context.Context, id int64) error
	UpdateExpires(ctx context.Context, id int64, expires_at *time.Time) error
}
type UrlShortenerService struct {
	repo    LinkRepository
	log     *slog.Logger
	hashing Hashing
}

func NewUrlShortenerService(log *slog.Logger, hash Hashing, repo LinkRepository) *UrlShortenerService {
	return &UrlShortenerService{repo: repo, log: log, hashing: hash}
}

func (s *UrlShortenerService) Shorten(ctx context.Context, sourceUrl string, expiresAt *time.Time) (string, error) {
	var linkId int64
	link, err := s.repo.GetLinkBySourceUrl(ctx, sourceUrl)
	if err != nil {
		if !errors.Is(err, reperrors.ErrLinkNotFound) {
			return "", err
		}

		// Проверяем срок действия только при создании ссылки.
		if expiresAt != nil && expiresAt.Before(time.Now().UTC()) {
			return "", serverrors.ErrExpired
		}
		linkId, err = s.repo.CreateLink(ctx, sourceUrl, expiresAt)
		if err != nil {
			return "", err
		}
	} else {
		linkId = link.ID

		// Существующую ссылку допускается сделать просроченной, проверка на expired исключена.
		if !reflecthelper.ComparePointers(link.ExpiresAt, expiresAt) {
			err = s.repo.UpdateExpires(ctx, link.ID, expiresAt)
			if err != nil {
				return "", err
			}
		}
	}

	return s.generateToken(linkId)
}

func (s *UrlShortenerService) Resolve(ctx context.Context, code string) (string, error) {
	const op = "service.urlShortenerService.Resolve"
	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	id, err := s.decodeToken(code)
	if err != nil {
		return "", err
	}
	if id == 0 {
		return "", errors.New("invalid code")
	}
	link, err := s.repo.GetLinkById(ctx, id)
	if err != nil {
		if errors.Is(err, reperrors.ErrLinkNotFound) {
			return "", serverrors.ErrURLNotFound
		}
		return "", err
	}
	err = s.repo.UpdateLastRequested(ctx, id)
	if err != nil {
		log.Error("failed to update last requested", slog.Any("error", err))
	}
	if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now().UTC()) {
		return "", serverrors.ErrURLExpired
	}
	return link.SourceUrl, nil
}

func (s *UrlShortenerService) generateToken(key int64) (string, error) {
	return s.hashing.EncodeInt64([]int64{key})
}

func (s *UrlShortenerService) decodeToken(token string) (int64, error) {
	decoded, err := s.hashing.DecodeInt64WithError(token)
	if err == nil && len(decoded) > 0 {
		return decoded[0], nil
	}
	return 0, err
}
