package domain_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alex8obrien/flow/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadyz_InitialState(t *testing.T) {
	r := domain.NewReadyz()
	exp := fmt.Sprintf(`{"version":"%s","overall":"down","db":"down"}`, domain.Version)

	// Check that the overall status is down
	assert.Equal(t, domain.DOWN, r.Overall())

	// Marshall the Readyz to JSON
	body, err := json.Marshal(r)
	require.NoError(t, err)
	assert.JSONEq(t, exp, string(body))
}

func TestReadyz_DatabaseChecked(t *testing.T) {
	r := domain.NewReadyz()
	r.DBChecked()
	exp := fmt.Sprintf(`{"version":"%s","overall":"ok","db":"ok"}`, domain.Version)

	// Check that the overall status is down
	assert.Equal(t, domain.OK, r.Overall())

	// Marshall the Readyz to JSON
	body, err := json.Marshal(r)
	require.NoError(t, err)
	assert.JSONEq(t, exp, string(body))
}
