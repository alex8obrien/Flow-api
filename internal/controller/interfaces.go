package controller

import (
	"context"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
)

type UserService interface {
	Register(ctx context.Context, name, email string) (*domain.User, error)
}
