package dto

import "time"

// LinkResponse reserved for future functionality
type LinkResponse struct {
	SourceUrl       string     `json:"source_url"`
	ExpiresAt       *time.Time `json:"expires_at"`
	LastRequestedAt *time.Time `json:"last_requested_at"`
}
