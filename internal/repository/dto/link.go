package dto

import "time"

// LinkDTO represents the link structure at the infrastructure level.
// Is used to transfer data between the database and business logic.
// Is redundant at this stage, as it coincides with the business entity,
// but is intended for future isolation and possible changes to the database schema.
type LinkDTO struct {
	ID              int64
	SourceUrl       string
	ExpiresAt       *time.Time
	LastRequestedAt *time.Time
	CreatedAt       time.Time
}
