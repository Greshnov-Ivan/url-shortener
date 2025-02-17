package dto

import "time"

type LinkDTO struct {
	ID              int64
	SourceUrl       string
	ExpiresAt       *time.Time
	LastRequestedAt *time.Time
	CreatedAt       time.Time
}
