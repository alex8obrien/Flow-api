package service_test

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/alex8obrien/flow/internal/domain"
	"github.com/alex8obrien/flow/internal/util"
)

type MockRepository struct {
}

type MockRepositoryError struct {
}

func NewMockRepository() MockRepository {
	return MockRepository{}
}

// FindByShortCode implements [service.LinkRepository].
func (m MockRepository) FindByShortCode(_ context.Context, shortCode string) (*domain.Link, error) {
	return &domain.Link{
		ID:             uuid.NewV7(),
		ShortCode:      shortCode,
		DestinationURL: "https://example.com",
		Clicks:         0,
		DisplayName:    "Example",
		Description:    "Example Description",
		CreatedAt:      time.Now(),
	}, nil
}

// Insert implements [service.LinkRepository].
func (m MockRepository) Insert(_ context.Context, link *domain.Link) (*domain.Link, error) {
	return &domain.Link{
		ID:             uuid.NewV7(),
		ShortCode:      util.RandomShortCode(),
		DestinationURL: link.DestinationURL,
		Clicks:         0,
		DisplayName:    link.DisplayName,
		Description:    link.Description,
		CreatedAt:      time.Now(),
	}, nil
}

// Ping implements [service.LinkRepository].
func (m MockRepository) Ping(_ context.Context) error {
	return nil
}

func NewMockRepositoryError() MockRepositoryError {
	return MockRepositoryError{}
}

// FindByShortCode implements [service.LinkRepository].
func (m MockRepositoryError) FindByShortCode(_ context.Context, shortCode string) (*domain.Link, error) {
	return nil, errors.New("failed to read from database")
}

// Insert implements [service.LinkRepository].
func (m MockRepositoryError) Insert(_ context.Context, link *domain.Link) (*domain.Link, error) {
	return nil, errors.New("failed to insert into database")
}

// Ping implements [service.LinkRepository].
func (m MockRepositoryError) Ping(ctx context.Context) error {
	return errors.New("failed to ping database")
}
