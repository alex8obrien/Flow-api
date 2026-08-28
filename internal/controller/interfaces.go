package controller

import (
	"context"
	"net/url"

	"github.com/alex8obrien/flow/internal/domain"
)

type LinkService interface {
	Create(ctx context.Context, destinationURL, displayName, description string) (*domain.Link, error)
	Redirect(ctx context.Context, shortCode string) (*url.URL, error)
}
