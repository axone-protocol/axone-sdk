package jwt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/axone-protocol/axone-sdk/auth"
)

// HTTPAuthHandler returns an HTTP handler that authenticates an auth.Identity and issues a related JWT token.
func (issuer *Issuer) HTTPAuthHandler(proxy auth.Proxy) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		credential, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, fmt.Errorf("failed to read request body credential: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id, err := proxy.Authenticate(context.Background(), credential)
		if err != nil {
			http.Error(writer, fmt.Errorf("failed to authenticate: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		token, err := issuer.IssueJWT(id)
		if err != nil {
			http.Error(writer, fmt.Errorf("failed to issue JWT: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write([]byte(token)); err != nil {
			http.Error(writer, fmt.Errorf("failed to write response: %w", err).Error(), http.StatusInternalServerError)
			return
		}
	})
}

// VerifyHTTPMiddleware returns an HTTP middleware that verifies the authenticity of a JWT token before forwarding the
// request to the next auth.AuthenticatedHandler providing the resolve auth.Identity.
func (issuer *Issuer) VerifyHTTPMiddleware(next auth.AuthenticatedHandler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := issuer.VerifyHTTPRequest(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		next(id, writer, request)
	})
}

// VerifyHTTPRequest checks the authenticity of the JWT token from the given HTTP request and returns the authenticated
// auth.Identity.
func (issuer *Issuer) VerifyHTTPRequest(r *http.Request) (*auth.Identity, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
		return nil, errors.New("couldn't find bearer token")
	}

	return issuer.VerifyJWT(authHeader[7:])
}
