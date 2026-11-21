package service

import (
	"context"

	"github.com/alex8obrien/flow/internal/domain"
)

type LinkRepository interface {
	Ping(ctx context.Context) error
	Insert(ctx context.Context, link *domain.Link) (*domain.Link, error)
	FindByShortCode(ctx context.Context, shortCode string) (*domain.Link, error)
}
