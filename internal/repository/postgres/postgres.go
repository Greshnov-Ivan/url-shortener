package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/entity"
	"url-shortener/internal/repository/reperrors"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CloseDB() error {
	return r.db.Close()
}

// реализовать интерфейс сканера type Scanner
func (r *Repository) CreateLink(ctx context.Context, sourceUrl string, expiresAt *time.Time) (int64, error) {
	query := `INSERT INTO links (source_url, expires_at, created_at) VALUES ($1, $2, $3) RETURNING id`
	nowUtc := time.Now().UTC()
	var id int64
	err := r.db.QueryRowContext(ctx, query, sourceUrl, expiresAt, nowUtc).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert link: %w", err)
	}
	return id, nil
}

func (r *Repository) GetLinkBySourceUrl(ctx context.Context, sourceUrl string) (*entity.Link, error) {
	query := "SELECT id, source_url, expires_at, created_at, last_requested_at FROM links WHERE source_url = $1"
	link := &entity.Link{}
	err := r.db.QueryRowContext(ctx, query, sourceUrl).Scan(&link.ID, &link.SourceUrl, &link.ExpiresAt, &link.CreatedAt, &link.LastRequestedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, reperrors.ErrLinkNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get link by source URL: %w", err)
	}
	return link, nil
}

func (r *Repository) GetLinkById(ctx context.Context, id int64) (*entity.Link, error) {
	query := "SELECT id, source_url, expires_at, created_at, last_requested_at FROM links WHERE id = $1"
	link := &entity.Link{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&link.ID, &link.SourceUrl, &link.ExpiresAt, &link.CreatedAt, &link.LastRequestedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, reperrors.ErrLinkNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get link by ID: %w", err)
	}
	return link, nil
}

func (r *Repository) UpdateLastRequested(ctx context.Context, id int64) error {
	nowUtc := time.Now().UTC()
	query := `UPDATE links SET last_requested_at = $2 WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id, nowUtc)
	if err != nil {
		return fmt.Errorf("failed to update last requested at: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return reperrors.ErrLinkNotFound
	}

	return nil
}

func (r *Repository) UpdateExpires(ctx context.Context, id int64, expiresAt *time.Time) error {
	query := `UPDATE links SET expires_at = $2 WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to update expiration: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return reperrors.ErrLinkNotFound
	}
	return nil
}
