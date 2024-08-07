package jwt

import "github.com/golang-jwt/jwt"

type ProxyClaims struct {
	jwt.StandardClaims
	Can Permissions `json:"can"`
}

type Permissions struct {
	Actions []string `json:"actions"`
}
