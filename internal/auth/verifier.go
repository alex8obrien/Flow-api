package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alex8obrien/flow/internal/util"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type LocalConfig struct {
	IssuerURL     string
	RawJWK        json.RawMessage
	SigningMethod string
	AudienceURL   string
}

type RemoteConfig struct {
	IssuerURL     string
	JwksURL       string
	SigningMethod string
	AudienceURL   string
}

type Verifier struct {
	issuer        string
	keys          keyfunc.Keyfunc
	signingMethod string
	audience      string
}

func NewLocalResolver(cfg LocalConfig) (*Verifier, error) {
	// Validate input
	if err := util.ValidateURL(cfg.IssuerURL); err != nil {
		return nil, fmt.Errorf("invalid issuer URL: %w", err)
	}
	if err := util.ValidateURL(cfg.AudienceURL); err != nil {
		return nil, fmt.Errorf("invalid audience URL: %w", err)
	}
	keys, err := keyfunc.NewJWKSetJSON(cfg.RawJWK)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK: %w", err)
	}

	return &Verifier{
		issuer:        cfg.IssuerURL,
		keys:          keys,
		signingMethod: cfg.SigningMethod,
		audience:      cfg.AudienceURL,
	}, nil
}

func NewRemoteVerifier(ctx context.Context, cfg RemoteConfig) (*Verifier, error) {
	// Validate input
	if err := util.ValidateURL(cfg.IssuerURL); err != nil {
		return nil, fmt.Errorf("invalid issuer URL: %w", err)
	}
	if err := util.ValidateURL(cfg.JwksURL); err != nil {
		return nil, fmt.Errorf("invalid JWKS URL: %w", err)
	}
	if err := util.ValidateURL(cfg.AudienceURL); err != nil {
		return nil, fmt.Errorf("invalid audience URL: %w", err)
	}
	keys, err := keyfunc.NewDefaultCtx(ctx, []string{cfg.JwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK: %w", err)
	}

	return &Verifier{
		issuer:        cfg.IssuerURL,
		keys:          keys,
		signingMethod: cfg.SigningMethod,
		audience:      cfg.AudienceURL,
	}, nil
}

func (v *Verifier) Verify(ctx context.Context, rawToken string) (*Claims, error) {
	// 0. Validate data input
	if rawToken == "" {
		return nil, ErrInvalidToken
	}
	claims := &Claims{}

	// 1. Parse the token
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		return v.keys.KeyfuncCtx(ctx)(token)
	},
		jwt.WithIssuer(v.issuer),
		jwt.WithValidMethods([]string{v.signingMethod}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithAudience(v.audience),
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. Check if the token is valid
	if !token.Valid || claims.Subject == "" || claims.IssuedAt == nil {
		return nil, ErrInvalidToken
	}

	// 3. Return claims
	return claims, nil
}
