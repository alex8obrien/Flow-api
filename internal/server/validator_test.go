package server_test

import (
	"net/http"
	"testing"

	"github.com/alex8obrien/flow/internal/server"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validationFixture struct {
	Name        string `validate:"required,min=3"`
	Destination string `validate:"required,url"`
}

func TestCustomValidator_ValidStruct(t *testing.T) {
	validator := server.NewCustomValidator()

	err := validator.Validate(validationFixture{
		Name:        "Example",
		Destination: "https://example.com",
	})

	require.NoError(t, err)
}

func TestCustomValidator_InvalidStruct(t *testing.T) {
	tests := []struct {
		name    string
		fixture validationFixture
	}{
		{name: "missing name", fixture: validationFixture{Destination: "https://example.com"}},
		{name: "name too short", fixture: validationFixture{Name: "x", Destination: "https://example.com"}},
		{name: "missing destination", fixture: validationFixture{Name: "Example"}},
		{name: "invalid destination", fixture: validationFixture{Name: "Example", Destination: "not-a-url"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := server.NewCustomValidator()

			err := validator.Validate(test.fixture)

			require.Error(t, err)
			httpError, ok := err.(*echo.HTTPError)
			require.True(t, ok)
			assert.Equal(t, http.StatusUnprocessableEntity, httpError.Code)
		})
	}
}
