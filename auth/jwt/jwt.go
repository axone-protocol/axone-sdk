// Package jwt brings a mean to manage JWT tokens on top of Axone network authentication mechanisms.
package jwt

import (
	"fmt"
	"time"

	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Issuer is an entity responsible to issue and verify JWT tokens.
type Issuer struct {
	secretKey []byte
	issuer    string
	ttl       time.Duration
}

// NewIssuer creates a new Issuer with the given secret key, issuer and token time-to-live.
func NewIssuer(secretKey []byte, issuer string, ttl time.Duration) *Issuer {
	return &Issuer{
		secretKey: secretKey,
		issuer:    issuer,
		ttl:       ttl,
	}
}

// IssueJWT forge and sign a new JWT token for the given authenticated auth.Identity.
func (issuer *Issuer) IssueJWT(identity *auth.Identity) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, ProxyClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  identity.DID,
			ExpiresAt: now.Add(issuer.ttl).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  now.Unix(),
			Issuer:    issuer.issuer,
			NotBefore: now.Unix(),
			Subject:   identity.DID,
		},
		Can: Permissions{
			Actions: identity.AuthorizedActions,
		},
	}).SignedString(issuer.secretKey)
}

// VerifyJWT checks the validity and the signature of the given JWT token and returns the authenticated identity.
func (issuer *Issuer) VerifyJWT(raw string) (*auth.Identity, error) {
	token, err := jwt.ParseWithClaims(raw, &ProxyClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return issuer.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*ProxyClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return &auth.Identity{
		DID:               claims.Subject,
		AuthorizedActions: claims.Can.Actions,
	}, nil
}
