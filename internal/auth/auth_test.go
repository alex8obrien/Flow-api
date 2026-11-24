package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/alex8obrien/flow/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	kid      = "test-key-1"
	audience = "https://api.example.com"
	issuer   = "https://auth.example.com"
)

// newRSAKey generates a new RSA key set.
func newRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return key
}

// newJwksJSON builds a JWKS document containing the public key.
func newJwksJSON(t *testing.T, pub *rsa.PublicKey, kid string) json.RawMessage {
	t.Helper()
	jwks := map[string]any{
		"keys": []map[string]string{{
			"kty": "RSA",
			"use": "sig",
			"alg": "RS256",
			"kid": kid,
			"n":   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		}},
	}
	b, err := json.Marshal(jwks)
	if err != nil {
		t.Fatalf("marshal jwks: %v", err)
	}
	return b
}

// signToken signs a JWT token and returns the complete string.
func signToken(t *testing.T, key *rsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = kid
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

// signTokenWithMethod signs a JWT token using the requested method.
func signTokenWithMethod(
	t *testing.T,
	key *rsa.PrivateKey,
	kid string,
	method jwt.SigningMethod,
	claims jwt.MapClaims,
) string {
	t.Helper()
	tok := jwt.NewWithClaims(method, claims)
	tok.Header["kid"] = kid
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

// accessClaims generates cognito-style access token claims.
func accessClaims(exp time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"sub":       "user-123",
		"iss":       issuer,
		"token_use": "access",
		"client_id": "test-client-id",
		"iat":       time.Now().Add(-time.Minute).Unix(),
		"exp":       exp.Unix(),
		"aud":       audience,
	}
}

func TestVerify(t *testing.T) {
	// Generate Keys & Local JWK
	goodKey := newRSAKey(t)
	otherKey := newRSAKey(t)
	jwks := newJwksJSON(t, &goodKey.PublicKey, kid)

	// Generate different sets of claims
	validClaims := accessClaims(time.Now().Add(time.Hour))

	expiredClaims := accessClaims(time.Now().Add(-time.Hour))

	wrongIssuerClaims := accessClaims(time.Now().Add(time.Hour))
	wrongIssuerClaims["iss"] = "https://wrong.example.com"

	wrongAudienceClaims := accessClaims(time.Now().Add(time.Hour))
	wrongAudienceClaims["aud"] = "https://wrong.example.com"

	missingSubClaims := accessClaims(time.Now().Add(time.Hour))
	delete(missingSubClaims, "sub")

	missingExpClaims := accessClaims(time.Now().Add(time.Hour))
	delete(missingExpClaims, "exp")

	missingIatClaims := accessClaims(time.Now().Add(time.Hour))
	delete(missingIatClaims, "iat")

	futureIatClaims := accessClaims(time.Now().Add(time.Hour))
	futureIatClaims["iat"] = time.Now().Add(time.Hour).Unix()

	malformedExpClaims := accessClaims(time.Now().Add(time.Hour))
	malformedExpClaims["exp"] = "not-a-timestamp"

	// Define different tests
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"empty token", "", true},
		{"malformed token", "not.a.jwt", true},
		{"valid", signToken(t, goodKey, kid, validClaims), false},
		{"expired", signToken(t, goodKey, kid, expiredClaims), true},
		{"wrong issuer", signToken(t, goodKey, kid, wrongIssuerClaims), true},
		{"wrong audience", signToken(t, goodKey, kid, wrongAudienceClaims), true},
		{"missing subject", signToken(t, goodKey, kid, missingSubClaims), true},
		{"missing expiration", signToken(t, goodKey, kid, missingExpClaims), true},
		{"missing issued at", signToken(t, goodKey, kid, missingIatClaims), true},
		{"future issued at", signToken(t, goodKey, kid, futureIatClaims), true},
		{"malformed expiration", signToken(t, goodKey, kid, malformedExpClaims), true},
		{"signed by wrong key", signToken(t, otherKey, kid, accessClaims(time.Now().Add(time.Hour))), true},
		{"unknown kid", signToken(t, goodKey, "nope", accessClaims(time.Now().Add(time.Hour))), true},
		{"missing kid", signToken(t, goodKey, "", accessClaims(time.Now().Add(time.Hour))), true},
		{
			"unsupported signing method",
			signTokenWithMethod(t, goodKey, kid, jwt.SigningMethodRS512, accessClaims(time.Now().Add(time.Hour))),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := auth.NewLocalResolver(auth.LocalConfig{
				IssuerURL:     issuer,
				RawJWK:        jwks,
				SigningMethod: jwt.SigningMethodRS256.Alg(),
				AudienceURL:   audience,
			})
			require.NoError(t, err)
			_, err = r.Verify(t.Context(), tt.token)
			if tt.wantErr {
				assert.ErrorIs(t, err, auth.ErrInvalidToken)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyReturnsClaims(t *testing.T) {
	key := newRSAKey(t)
	jwks := newJwksJSON(t, &key.PublicKey, kid)
	r, err := auth.NewLocalResolver(auth.LocalConfig{
		IssuerURL:     issuer,
		RawJWK:        jwks,
		SigningMethod: jwt.SigningMethodRS256.Alg(),
		AudienceURL:   audience,
	})
	require.NoError(t, err)

	token := signToken(t, key, kid, accessClaims(time.Now().Add(time.Hour)))
	claims, err := r.Verify(t.Context(), token)

	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, "user-123", claims.Subject)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Equal(t, "access", claims.TokenUse)
	assert.Equal(t, "test-client-id", claims.ClientID)
	assert.Contains(t, claims.Audience, audience)
}
