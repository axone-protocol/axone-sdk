package jwt

import (
	"fmt"
	"github.com/axone-protocol/axone-sdk/auth"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Factory struct {
	secretKey []byte
	issuer    string
	ttl       time.Duration
}

func NewFactory(secretKey []byte, issuer string, ttl time.Duration) *Factory {
	return &Factory{
		secretKey: secretKey,
		issuer:    issuer,
		ttl:       ttl,
	}
}

func (f *Factory) IssueJWT(identity *auth.Identity) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, ProxyClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  identity.DID,
			ExpiresAt: now.Add(f.ttl).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  now.Unix(),
			Issuer:    f.issuer,
			NotBefore: now.Unix(),
			Subject:   identity.DID,
		},
		Can: Permissions{
			Actions: identity.AuthorizedActions,
		},
	}).SignedString(f.secretKey)
}

func (f *Factory) VerifyJWT(raw string) (*auth.Identity, error) {
	token, err := jwt.ParseWithClaims(raw, &ProxyClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return f.secretKey, nil
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
