package db_test

import (
	"testing"

	"github.com/alex8obrien/flow/internal/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPool(t *testing.T) {
	pool, err := db.NewPool(t.Context(), "postgresql://example.com:5432/postgres?user=postgres&password=postgres")
	require.NoError(t, err)
	require.NotNil(t, pool)
	t.Cleanup(pool.Close)
}

func TestNewPool_InvalidDsn(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{name: "Empty DSN", url: ""},
		{name: "Malformed URL", url: "://postgresql"},
		{name: "Invalid Protocol Only", url: "psql://"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool, err := db.NewPool(t.Context(), test.url)
			require.Error(t, err)
			assert.Nil(t, pool)
		})
	}
}
