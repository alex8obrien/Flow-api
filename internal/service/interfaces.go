package service

import (
	"context"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
)

type UserRepository interface {
	// FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Insert(ctx context.Context, user *domain.User) (*domain.User, error)
}
