package environment_test

import (
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

func clearEnvironment(t *testing.T) {
	t.Helper()
	for _, key := range requiredEnvironmentKeys {
		value, exists := os.LookupEnv(key)
		require.NoError(t, os.Unsetenv(key))
		t.Cleanup(func() {
			if exists {
				require.NoError(t, os.Setenv(key, value))
			} else {
				require.NoError(t, os.Unsetenv(key))
			}
		})
	}
}

func writeEnvironmentFile(t *testing.T, omit string) string {
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

	var contents string
	for key, value := range values {
		if key != omit {
			contents += key + "=" + value + "\n"
		}
	}

	path := filepath.Join(t.TempDir(), ".env")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
	return path
}

func TestInitEnv_LoadsCompleteFile(t *testing.T) {
	clearEnvironment(t)
	path := writeEnvironmentFile(t, "")

	require.NoError(t, environment.InitEnv(path))
	assert.Equal(t, ":8080", environment.GetEnv(environment.Address))
	assert.Equal(t, "flow", environment.GetEnv(environment.ApplicationName))
	assert.Equal(t, "RS256", environment.GetEnv(environment.JwtSigningMethod))
	assert.Equal(t, "https://api.example.com", environment.GetEnv(environment.JwtAudience))
}

func TestInitEnv_MissingVariables(t *testing.T) {
	for _, key := range requiredEnvironmentKeys {
		t.Run(key, func(t *testing.T) {
			clearEnvironment(t)
			path := writeEnvironmentFile(t, key)

			err := environment.InitEnv(path)

			require.Error(t, err)
			assert.Contains(t, err.Error(), key)
		})
	}
}

func TestInitEnv_MissingFile(t *testing.T) {
	clearEnvironment(t)

	err := environment.InitEnv(filepath.Join(t.TempDir(), "missing.env"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load environment file")
}

func TestGetEnv(t *testing.T) {
	t.Setenv("FLOW_TEST_VALUE", "test-value")

	assert.Equal(t, "test-value", environment.GetEnv("FLOW_TEST_VALUE"))
}
