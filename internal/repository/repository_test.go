package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"
	"uuid"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
	"github.com/alex8obrien/echo-v5-template/internal/repository"
	"github.com/alex8obrien/echo-v5-template/internal/test/telemetry"
	"github.com/alex8obrien/echo-v5-template/internal/util"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Pool pgxmock.PgxPoolIface

func TestMain(m *testing.M) {
	var err error
	if Pool, err = pgxmock.NewPool(); err != nil {
		return
	}
	defer Pool.Close()
	if code := m.Run(); code != 0 {
		return
	}
}

func TestUserRepository_Insert(t *testing.T) {
	expectedID := uuid.NewV7()
	query := "INSERT INTO links (short_code, destination_url, display_name, description) VALUES ($1, $2, $3, $4) RETURNING id, clicks, created_at;"
	shortCode, err := util.GenerateTruncatedShortCode("https://google.com", domain.ShortCodeMaximumLength)
	require.NoError(t, err)

	Pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(shortCode, "https://google.com", "Google", "").
		WillReturnRows(pgxmock.NewRows([]string{"id", "clicks", "created_at"}).AddRow(expectedID, 0, time.Now()))

	repo := repository.NewLinkRepository(Pool, telemetry.NewTracer())
	link, err := repo.Insert(t.Context(), &domain.Link{
		ShortCode:      shortCode,
		DestinationURL: "https://google.com",
		DisplayName:    "Google",
		Description:    "",
	})

	require.NoError(t, err)
	assert.Equal(t, expectedID, link.ID)
	require.NoError(t, Pool.ExpectationsWereMet())
}

func TestUserRepository_Insert_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	query := "INSERT INTO links (short_code, destination_url, display_name, description) VALUES ($1, $2, $3, $4) RETURNING id, clicks, created_at;"
	shortCode, err := util.GenerateTruncatedShortCode("https://google.com", domain.ShortCodeMaximumLength)
	require.NoError(t, err)

	Pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(shortCode, "https://google.com", "Google", "").
		WillReturnError(context.Canceled)

	repo := repository.NewLinkRepository(Pool, telemetry.NewTracer())
	link, err := repo.Insert(ctx, &domain.Link{
		ShortCode:      shortCode,
		DestinationURL: "https://google.com",
		DisplayName:    "Google",
		Description:    "",
	})

	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, link)
	require.Error(t, err)
	require.NoError(t, Pool.ExpectationsWereMet())
}
