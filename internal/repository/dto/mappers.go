package dto

import (
	"time"
	"url-shortener/internal/entity"
)

// MapLinkToDTO reserved for future functionality
func MapLinkToDTO(link *entity.Link) *LinkDTO {
	if link == nil {
		return nil
	}
	return &LinkDTO{
		ID:              link.ID,
		SourceUrl:       link.SourceUrl,
		ExpiresAt:       link.ExpiresAt,
		LastRequestedAt: link.LastRequestedAt,
		CreatedAt:       time.Now().UTC(),
	}
}

func MapDTOToLink(dto *LinkDTO) *entity.Link {
	if dto == nil {
		return nil
	}
	return &entity.Link{
		ID:              dto.ID,
		SourceUrl:       dto.SourceUrl,
		ExpiresAt:       dto.ExpiresAt,
		LastRequestedAt: dto.LastRequestedAt,
	}
}
