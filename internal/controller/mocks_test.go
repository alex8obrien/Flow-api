package controller_test

import (
	"context"
	"uuid"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
)

type MockService struct{}

func NewMockService() MockService {
	return MockService{}
}

func (srv MockService) Register(_ context.Context, name, email string) (*domain.User, error) {
	return &domain.User{
		ID:    uuid.NewV7(),
		Name:  name,
		Email: email,
	}, nil
}
