package jwt

import (
	"context"
	"errors"
	"github.com/axone-protocol/axone-sdk/auth"
	"io"
	"net/http"
)

func (f *Factory) HTTPAuthHandler(proxy auth.Proxy) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		credential, err := io.ReadAll(request.Body)
		if err != nil {
			// ...
		}

		id, err := proxy.Authenticate(context.Background(), credential)
		if err != nil {
			// ...
		}

		token, err := f.IssueJWT(id)
		if err != nil {
			// ...
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write([]byte(token)); err != nil {
			// ...
		}
	})
}

func (f *Factory) VerifyHTTPMiddleware(next auth.AuthenticatedHandler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := f.VerifyHTTPRequest(request)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(err.Error()))
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
