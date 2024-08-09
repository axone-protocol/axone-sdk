package storage

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/auth/jwt"
	axonehttp "github.com/axone-protocol/axone-sdk/http"
	"github.com/gorilla/mux"
)

func (p *Proxy) HTTPConfigurator(jwtSecretKey []byte, jwtTTL time.Duration) axonehttp.Option {
	jwtFactory := jwt.NewFactory(jwtSecretKey, p.key.DID, jwtTTL)

	return axonehttp.WithOptions(
		axonehttp.WithRoute(http.MethodGet, "/authenticate", jwtFactory.HTTPAuthHandler(p)),
		axonehttp.WithRoute(http.MethodGet, "/{path}", jwtFactory.VerifyHTTPMiddleware(p.HTTPReadHandler())),
		axonehttp.WithRoute(http.MethodPost, "/{path}", jwtFactory.VerifyHTTPMiddleware(p.HTTPStoreHandler())),
	)
}

func (p *Proxy) HTTPReadHandler() auth.AuthenticatedHandler {
	return func(id *auth.Identity, writer http.ResponseWriter, request *http.Request) {
		resource, err := p.Read(context.Background(), id, mux.Vars(request)["path"])
		if err != nil {
			// ...
			return
		}

		writer.WriteHeader(http.StatusOK)
		io.Copy(writer, resource)
	}
}

func (p *Proxy) HTTPStoreHandler() auth.AuthenticatedHandler {
	return func(id *auth.Identity, writer http.ResponseWriter, request *http.Request) {
		vc, err := p.Store(context.Background(), id, mux.Vars(request)["path"], request.Body)
		if err != nil {
			// ...
			return
		}

		writer.WriteHeader(http.StatusOK)
		io.Copy(writer, vc)
	}
}
