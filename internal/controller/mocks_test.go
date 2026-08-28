package controller_test

import (
	"context"
	"net/url"
	"time"
	"uuid"

	"github.com/alex8obrien/flow/internal/domain"
	"github.com/alex8obrien/flow/internal/util"
)

type MockService struct{}

func NewMockService() MockService {
	return MockService{}
}

func (srv MockService) Create(
	_ context.Context,
	_, _, _ string,
) (*domain.Link, error) {
	return &domain.Link{
		ID:             uuid.NewV7(),
		ShortCode:      util.RandomShortCode(),
		DestinationURL: "https://example.com",
		Clicks:         0,
		DisplayName:    "Example",
		CreatedAt:      time.Now(),
	}, nil
}

func (srv MockService) Redirect(_ context.Context, _ string) (*url.URL, error) {
	return url.Parse("https://example.com")
}
