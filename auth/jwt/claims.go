package jwt

import "github.com/golang-jwt/jwt"

// ProxyClaims is the set of claims that are included in the JWT token.
type ProxyClaims struct {
	jwt.StandardClaims
	Can Permissions `json:"can"`
}

// Permissions is the set of permissions that are included in the JWT token.
type Permissions struct {
	Actions []string `json:"actions"`
}
