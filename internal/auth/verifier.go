package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/alex8obrien/flow/internal/util"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

var errInvalidToken = errors.New("invalid token")

type Config struct {
	IssuerURL     string
	JwksURL       string
	SigningMethod string
}

type Verifier struct {
	issuer     string
	keys       keyfunc.Keyfunc
	signingKey string
}

func NewVerifier(ctx context.Context, cfg Config) (*Verifier, error) {
	// Validate input
	if err := util.ValidateURL(cfg.IssuerURL); err != nil {
		return nil, fmt.Errorf("invalid issuer URL: %w", err)
	}
	if err := util.ValidateURL(cfg.JwksURL); err != nil {
		return nil, fmt.Errorf("invalid JWKS URL: %w", err)
	}
	keys, err := keyfunc.NewDefaultCtx(ctx, []string{cfg.JwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK: %w", err)
	}

	return &Verifier{
		issuer:     cfg.IssuerURL,
		keys:       keys,
		signingKey: cfg.SigningMethod,
	}, nil
}

func (v *Verifier) Verify(ctx context.Context, rawToken string) (*Claims, error) {
	// 0. Validate data input
	if rawToken == "" {
		return nil, errInvalidToken
	}
	claims := &Claims{}

	// 1. Parse the token
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		return v.keys.KeyfuncCtx(ctx)(token)
	},
		jwt.WithIssuer(v.issuer),
		jwt.WithValidMethods([]string{v.signingKey}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
	if err != nil {
		return nil, errInvalidToken
	}

	// 2. Check if the token is valid
	if !token.Valid || claims.Subject == "" {
		return nil, errInvalidToken
	}

	// 3. Return claims
	return claims, nil
}
