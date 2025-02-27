package entity

import "time"

// Link represents the core business entity for a shortened link.
// It contains essential information about the original URL and its lifecycle.
// This entity is used in the business logic layer and remains independent
// of database-specific implementations.
type Link struct {
	// ID is the unique identifier of the link.
	ID int64

	// SourceUrl is the original URL that is being shortened.
	SourceUrl string

	// ExpiresAt is an optional expiration date after which the link becomes invalid.
	ExpiresAt *time.Time

	// LastRequestedAt is the timestamp of the last time the link was accessed.
	LastRequestedAt *time.Time
}
