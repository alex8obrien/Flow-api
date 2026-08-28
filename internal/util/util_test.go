package util_test

import (
	"testing"

	"github.com/alex8obrien/echo-v5-template/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtil_CanonicaliseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Standard URL",
			"https://google.com",
			"https://google.com",
		},
		{
			"Trailing Slash",
			"https://google.com/",
			"https://google.com",
		},
		{
			"Removes Frag",
			"https://en.wikipedia.org/wiki/Main_Page#content",
			"https://en.wikipedia.org/wiki/Main_Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := util.CanonicaliseURL(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestUtil_GenerateShortCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Google",
			"https://google.com",
			"1BlRRSRMNECq2ZS4JeLoum3uucxjLDREag7BkRr8sGb",
		},
		{
			"Amazon",
			"https://amazon.com",
			"odCsksaz581Vz3QaL9FfMVKtiTpjg99JI7f3uj7s9Sz",
		},
		{
			"alex8obrien",
			"https://alex8obrien.com",
			"xx85qSHTwVjvohfIl1UDDyEwfGBfnto65XQz3uiUWes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := util.GenerateShortCode(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}
}
