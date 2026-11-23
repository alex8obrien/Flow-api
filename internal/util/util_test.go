package util_test

import (
	"testing"

	"github.com/alex8obrien/flow/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanonicaliseURL(t *testing.T) {
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

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Google",
			"https://google.com",
			"NOJoH5uzNXHfrpIr8gVSELaoHbVXWt7ZRkRfgsJvRQB",
		},
		{
			"Amazon",
			"https://amazon.com",
			"JzIc60Uk5GpcCvxqDVtorhqOrv3Ud5HwilpmqHkWFWN",
		},
		{
			"alex8obrien",
			"https://alex8obrien.com",
			"21PWpJlA9Ul58CtpfDyzLDQiL34eYIGPjyCJvI5pRNfB",
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
