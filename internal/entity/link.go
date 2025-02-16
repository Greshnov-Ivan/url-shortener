package entity

import "time"

type Link struct {
	ID              int64      `json:"id"`
	SourceUrl       string     `json:"source_url"`
	ExpiresAt       *time.Time `json:"expires_at"`
	LastRequestedAt *time.Time `json:"last_requested_at"`
	CreatedAt       time.Time  `json:"created_at"`
}
