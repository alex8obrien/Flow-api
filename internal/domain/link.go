package domain

import (
	"errors"
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

func NewLink(destinationURL, displayName, description string) (*Link, error) {
	// Validate destination URL
	url, err := util.CanonicaliseURL(destinationURL)
	if err != nil {
		return nil, err
	}

	// Generate short code
	code, err := util.GenerateTruncatedShortCode(destinationURL, ShortCodeMaximumLength)
	if err != nil {
		return nil, err
	}

	return &Link{
		ShortCode:      code,
		DestinationURL: url,
		DisplayName:    displayName,
		Description:    description,
	}, nil
}

var (
	ErrInvalidDestinationURL = errors.New("invalid destination url")
	ErrInvalidDisplayName    = errors.New("invalid display name")
)

const ShortCodeMaximumLength = 8
