package repository_test

import (
	"context"
	"regexp"
	"testing"
	"uuid"

	"github.com/alex8obrien/echo-v5-template/internal/domain"
	"github.com/alex8obrien/echo-v5-template/internal/repository"
	"github.com/alex8obrien/echo-v5-template/internal/test/telemetry"

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
	query := "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id"
	Pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("alice", "alice@example.com").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

	repo := repository.NewUserRepository(Pool, telemetry.NewTracer())
	user, err := repo.Insert(t.Context(), &domain.User{
		Name:  "alice",
		Email: "alice@example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, expectedID, user.ID)
	require.NoError(t, Pool.ExpectationsWereMet())
}

func TestUserRepository_Insert_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	query := "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id"
	Pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("alice", "alice@example.com").
		WillReturnError(context.Canceled)

	repo := repository.NewUserRepository(Pool, telemetry.NewTracer())

	result, err := repo.Insert(ctx, &domain.User{
		Name:  "alice",
		Email: "alice@example.com",
	})

	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, result)
	require.Error(t, err)
}
