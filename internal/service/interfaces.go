package service

import (
	"context"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
)

type LinkRepository interface {
	Insert(ctx context.Context, link *domain.Link) (*domain.Link, error)
	FindByShortCode(ctx context.Context, shortCode string) (*domain.Link, error)
}
