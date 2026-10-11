package environment

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	Address          = "ADDRESS"
	CertFile         = "CERT_FILE"
	KeyFile          = "KEY_FILE"
	AllowOrigins     = "ALLOW_ORIGIN"
	Dsn              = "POSTGRES_DSN"
	ApplicationName  = "APPLICATION_NAME"
	JwtIssuerURL     = "JWT_ISSUER_URL"
	JwksURL          = "JWKS_URL"
	JwtSigningMethod = "JWT_SIGNING_METHOD"
)

const errMessage = "missing environment variable %s"

// InitEnv loads environment file at the specified path.
// Confirms that every required variable has a set value.
func InitEnv(path string) error {
	if err := godotenv.Overload(path); err != nil {
		return fmt.Errorf("failed to load environment file. path: %s. error: %w", path, err)
	}

	if os.Getenv(Address) == "" {
		return fmt.Errorf(errMessage, Address)
	}
	if os.Getenv(CertFile) == "" {
		return fmt.Errorf(errMessage, CertFile)
	}
	if os.Getenv(KeyFile) == "" {
		return fmt.Errorf(errMessage, KeyFile)
	}
	if os.Getenv(AllowOrigins) == "" {
		return fmt.Errorf(errMessage, AllowOrigins)
	}
	if os.Getenv(Dsn) == "" {
		return fmt.Errorf(errMessage, Dsn)
	}
	if os.Getenv(ApplicationName) == "" {
		return fmt.Errorf(errMessage, ApplicationName)
	}
	if os.Getenv(JwtIssuerURL) == "" {
		return fmt.Errorf(errMessage, JwtIssuerURL)
	}
	if os.Getenv(JwksURL) == "" {
		return fmt.Errorf(errMessage, JwksURL)
	}
	if os.Getenv(JwtSigningMethod) == "" {
		return fmt.Errorf(errMessage, JwtSigningMethod)
	}

	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
