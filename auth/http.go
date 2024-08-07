package auth

import "net/http"

type AuthenticatedHandler func(*Identity, http.ResponseWriter, *http.Request)
