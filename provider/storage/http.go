package storage

import (
	"context"
	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/auth/jwt"
	axonehttp "github.com/axone-protocol/axone-sdk/http"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

func (s *Proxy) HTTPConfigurator(jwtSecretKey []byte, jwtTTL time.Duration) axonehttp.Option {
	jwtFactory := jwt.NewFactory(jwtSecretKey, s.key.DID, jwtTTL)

	return axonehttp.WithOptions(
		axonehttp.WithRoute(http.MethodGet, "/authenticate", jwtFactory.HTTPAuthHandler(s)),
		axonehttp.WithRoute(http.MethodGet, "/{path}}", jwtFactory.VerifyHTTPMiddleware(s.HTTPReadHandler())),
		axonehttp.WithRoute(http.MethodPost, "/{path}}", jwtFactory.VerifyHTTPMiddleware(s.HTTPStoreHandler())),
	)
}

func (s *Proxy) HTTPReadHandler() auth.AuthenticatedHandler {
	return func(id *auth.Identity, writer http.ResponseWriter, request *http.Request) {
		resource, err := s.Read(context.Background(), id, mux.Vars(request)["path"])
		if err != nil {
			// ...
			return
		}

		writer.WriteHeader(http.StatusOK)
		io.Copy(writer, resource)
	}
}

func (s *Proxy) HTTPStoreHandler() auth.AuthenticatedHandler {
	return func(id *auth.Identity, writer http.ResponseWriter, request *http.Request) {
		vc, err := s.Store(context.Background(), id, mux.Vars(request)["path"], request.Body)
		if err != nil {
			// ...
			return
		}

		writer.WriteHeader(http.StatusOK)
		io.Copy(writer, vc)
	}
}
