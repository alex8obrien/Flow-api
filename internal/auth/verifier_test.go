package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func rsaJWTTestFixture(t *testing.T) (*rsa.PrivateKey, keyfunc.Keyfunc) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKey := &privateKey.PublicKey
	jwks := map[string]any{
		"keys": []map[string]any{{
			"kty": "RSA",
			"kid": "test-key",
			"use": "sig",
			"alg": "RS256",
			"n":   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
		}},
	}
	raw, err := json.Marshal(jwks)
	require.NoError(t, err)

	keys, err := keyfunc.NewJWKSetJSON(raw)
	require.NoError(t, err)
	return privateKey, keys
}

func signedJWT(t *testing.T, privateKey *rsa.PrivateKey, claims *Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key"
	raw, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return raw
}

func validClaims() *Claims {
	now := time.Now()
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://auth.example.com/",
			Subject:   "user-123",
			Audience:  jwt.ClaimStrings{"https://api.example.com"},
			IssuedAt:  jwt.NewNumericDate(now.Add(-time.Minute)),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
		TokenUse: "access",
		ClientID: "flow-api",
	}
}

func TestVerifier_ValidToken(t *testing.T) {
	privateKey, keys := rsaJWTTestFixture(t)
	verifier := &Verifier{
		issuer:        "https://auth.example.com/",
		keys:          keys,
		signingMethod: jwt.SigningMethodRS256.Alg(),
		audience:      "https://api.example.com",
	}

	claims, err := verifier.Verify(
		t.Context(),
		signedJWT(t, privateKey, validClaims()),
	)

	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.Subject)
	assert.Equal(t, "access", claims.TokenUse)
}

func TestVerifier_InvalidTokens(t *testing.T) {
	privateKey, keys := rsaJWTTestFixture(t)
	verifier := &Verifier{
		issuer:        "https://auth.example.com/",
		keys:          keys,
		signingMethod: jwt.SigningMethodRS256.Alg(),
		audience:      "https://api.example.com",
	}

	tests := []struct {
		name  string
		token func(*rsa.PrivateKey) string
		raw   string
	}{
		{name: "empty", raw: ""},
		{
			name: "wrong issuer",
			token: func(key *rsa.PrivateKey) string {
				claims := validClaims()
				claims.Issuer = "https://wrong.example.com/"
				return signedJWT(t, key, claims)
			},
		},
		{
			name: "wrong audience",
			token: func(key *rsa.PrivateKey) string {
				claims := validClaims()
				claims.Audience = jwt.ClaimStrings{"https://wrong.example.com"}
				return signedJWT(t, key, claims)
			},
		},
		{
			name: "expired",
			token: func(key *rsa.PrivateKey) string {
				claims := validClaims()
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute))
				return signedJWT(t, key, claims)
			},
		},
		{
			name: "missing subject",
			token: func(key *rsa.PrivateKey) string {
				claims := validClaims()
				claims.Subject = ""
				return signedJWT(t, key, claims)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw := test.raw
			if test.token != nil {
				raw = test.token(privateKey)
			}
			claims, err := verifier.Verify(t.Context(), raw)

			require.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestNewVerifier_InvalidConfiguration(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "invalid issuer",
			cfg: Config{
				IssuerURL:     "://bad",
				JwksURL:       "https://auth.example.com/.well-known/jwks.json",
				AudienceURL:   "https://api.example.com",
				SigningMethod: "RS256",
			},
		},
		{
			name: "invalid JWKS URL",
			cfg: Config{
				IssuerURL:     "https://auth.example.com/",
				JwksURL:       "://bad",
				AudienceURL:   "https://api.example.com",
				SigningMethod: "RS256",
			},
		},
		{
			name: "invalid audience",
			cfg: Config{
				IssuerURL:     "https://auth.example.com/",
				JwksURL:       "https://auth.example.com/.well-known/jwks.json",
				AudienceURL:   "://bad",
				SigningMethod: "RS256",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier, err := NewVerifier(t.Context(), test.cfg)

			require.Error(t, err)
			assert.Nil(t, verifier)
		})
	}
}
