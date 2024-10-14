package auth

import "net/http"

// AuthenticatedHandler is a handler that requires an authenticated user.
//
// It is intended to be used in combination with a middleware handler that verifies the user's identity, for example
// with the jwt middleware:
//
//	jwtFactory.VerifyHTTPMiddleware(func(id *auth.Identity, w http.ResponseWriter, r *http.Request) {
//	  // Your handler logic here
//	})
type AuthenticatedHandler func(*Identity, http.ResponseWriter, *http.Request)
