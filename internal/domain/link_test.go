package domain_test

import (
	"testing"

	"github.com/alex8obrien/flow/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLink(t *testing.T) {
	link, err := domain.NewLink(
		"https://google.com/",
		"Google",
		"Search",
	)

	require.NoError(t, err)
	assert.Equal(t, "https://google.com", link.DestinationURL)
	assert.Equal(t, "Google", link.DisplayName)
	assert.Equal(t, "Search", link.Description)
	assert.Len(t, link.ShortCode, domain.ShortCodeMaximumLength)
	assert.Zero(t, link.Clicks)
	assert.True(t, link.CreatedAt.IsZero())
}

func TestNewLink_InvalidDestinationURLs(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{name: "empty", url: ""},
		{name: "missing host", url: "https://"},
		{name: "HTTP is not allowed", url: "http://example.com"},
		{name: "relative URL", url: "/example"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			link, err := domain.NewLink(test.url, "Example", "")

			require.Error(t, err)
			assert.Nil(t, link)
		})
	}
}
