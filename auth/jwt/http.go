package jwt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/axone-protocol/axone-sdk/auth"
)

func (f *Factory) HTTPAuthHandler(proxy auth.Proxy) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		credential, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, fmt.Errorf("failed to read request body credential: %w", err).Error(), http.StatusBadRequest)
			return
		}

		id, err := proxy.Authenticate(context.Background(), credential)
		if err != nil {
			http.Error(writer, fmt.Errorf("failed to authenticate: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		token, err := f.IssueJWT(id)
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

func (f *Factory) VerifyHTTPMiddleware(next auth.AuthenticatedHandler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := f.VerifyHTTPRequest(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		next(id, writer, request)
	})
}

func (f *Factory) VerifyHTTPRequest(r *http.Request) (*auth.Identity, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
		return nil, errors.New("couldn't find bearer token")
	}

	return f.VerifyJWT(authHeader[7:])
}
