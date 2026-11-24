package domain

import (
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/alex8obrien/flow/internal/util"
)

type Link struct {
	ID             uuid.UUID
	ShortCode      string
	DestinationURL string
	Clicks         int
	DisplayName    string
	Description    string
	CreatedAt      time.Time
	DeletedAt      time.Time
}

const ShortCodeMaximumLength = 8

var (
	ErrInvalidDestinationURL = errors.New("invalid destination url")
	ErrInvalidDisplayName    = errors.New("invalid display name")
)

func NewLink(destinationURL, displayName, description string) (*Link, error) {
	// Validate destination URL
	url, err := util.CanonicaliseURL(destinationURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidDestinationURL, err)
	}

	// Check for display name
	if displayName == "" {
		return nil, ErrInvalidDisplayName
	}

	// Generate short code
	code, err := util.GenerateTruncatedShortCode(destinationURL, ShortCodeMaximumLength)
	if err != nil {
		return nil, err
	}

	return &Link{
		ID:             uuid.Nil(),
		ShortCode:      code,
		DestinationURL: url,
		DisplayName:    displayName,
		Description:    description,
	}, nil
}
