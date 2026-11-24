package environment_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alex8obrien/flow/internal/environment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var requiredEnvironmentKeys = []string{
	environment.Address,
	environment.CertFile,
	environment.KeyFile,
	environment.AllowOrigins,
	environment.Dsn,
	environment.ApplicationName,
	environment.JwtIssuerURL,
	environment.JwksURL,
	environment.JwtSigningMethod,
	environment.JwtAudience,
}

func createEnvFile(t *testing.T, omit string) string {
	t.Helper()
	values := map[string]string{
		environment.Address:          ":8080",
		environment.CertFile:         "/cert.pem",
		environment.KeyFile:          "/key.pem",
		environment.AllowOrigins:     "https://example.com",
		environment.Dsn:              "postgres://app:password@localhost/app",
		environment.ApplicationName:  "flow",
		environment.JwtIssuerURL:     "https://auth.example.com/",
		environment.JwksURL:          "https://auth.example.com/.well-known/jwks.json",
		environment.JwtSigningMethod: "RS256",
		environment.JwtAudience:      "https://api.example.com",
	}

	// Loop through values and write to file
	var buf bytes.Buffer
	for key, value := range values {
		if key != omit {
			_, err := fmt.Fprintf(&buf, "%s=%s\n", key, value)
			require.NoError(t, err)
		}
	}

	path := filepath.Join(t.TempDir(), ".env")
	err := os.WriteFile(path, buf.Bytes(), 0o600)
	require.NoError(t, err)
	return path
}

func TestInitEnv_LoadsCompleteFile(t *testing.T) {
	path := createEnvFile(t, "")

	require.NoError(t, environment.InitEnv(path))
	assert.Equal(t, ":8080", environment.GetEnv(environment.Address))
	assert.Equal(t, "/cert.pem", environment.GetEnv(environment.CertFile))
	assert.Equal(t, "/key.pem", environment.GetEnv(environment.KeyFile))
	assert.Equal(t, "https://example.com", environment.GetEnv(environment.AllowOrigins))
	assert.Equal(t, "postgres://app:password@localhost/app", environment.GetEnv(environment.Dsn))
	assert.Equal(t, "flow", environment.GetEnv(environment.ApplicationName))
	assert.Equal(t, "https://auth.example.com/", environment.GetEnv(environment.JwtIssuerURL))
	assert.Equal(t, "https://auth.example.com/.well-known/jwks.json", environment.GetEnv(environment.JwksURL))
	assert.Equal(t, "RS256", environment.GetEnv(environment.JwtSigningMethod))
	assert.Equal(t, "https://api.example.com", environment.GetEnv(environment.JwtAudience))
}

func TestInitEnv_MissingVariables(t *testing.T) {
	for _, key := range requiredEnvironmentKeys {
		t.Run(key, func(t *testing.T) {
			path := createEnvFile(t, key)
			t.Setenv(key, "")
			err := environment.InitEnv(path)

			require.Error(t, err)
			assert.Contains(t, err.Error(), key)
		})
	}
}

func TestInitEnv_MissingFile(t *testing.T) {
	err := environment.InitEnv(filepath.Join(t.TempDir(), "missing.env"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load environment file")
}

func TestGetEnv(t *testing.T) {
	t.Setenv("FLOW_TEST_VALUE", "test-value")

	assert.Equal(t, "test-value", environment.GetEnv("FLOW_TEST_VALUE"))
}
