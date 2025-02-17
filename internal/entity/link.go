package entity

import "time"

type Link struct {
	ID              int64
	SourceUrl       string
	ExpiresAt       *time.Time
	LastRequestedAt *time.Time
}
