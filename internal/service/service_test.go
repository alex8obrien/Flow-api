package service_test

import (
	"os"
	"testing"

	"github.com/alex8obrien/flow/internal/service"
	"github.com/alex8obrien/flow/internal/test/telemetry"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var srv service.LinkService
var srvErr service.LinkService

func TestMain(m *testing.M) {
	srv = service.NewLinkService(NewMockRepository(), telemetry.NewTracer())
	srvErr = service.NewLinkService(NewMockRepositoryError(), telemetry.NewTracer())
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	link, err := srv.Create(t.Context(), "https://example.com", "Example", "Example Description")
	require.NoError(t, err)
	assert.NotNil(t, link)
}

func TestCreate_InvalidData(t *testing.T) {
	tests := []struct {
		name           string
		destinationURL string
		displayName    string
		description    string
	}{
		{
			name:           "Missing Destination URL",
			destinationURL: "",
			displayName:    "Example",
			description:    "Example Description",
		},
		{
			name:           "Missing Display Name",
			destinationURL: "https://example.com",
			displayName:    "",
			description:    "Example Description",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			link, err := srv.Create(t.Context(), test.destinationURL, test.displayName, test.description)
			require.Error(t, err)
			assert.Nil(t, link)
		})
	}
}

func TestCreate_InsertionFailure(t *testing.T) {
	link, err := srvErr.Create(t.Context(), "https://example.com", "Example", "Example Description")
	require.Error(t, err)
	assert.Nil(t, link)
}
