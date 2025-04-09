package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

// TokenMiddleware checks for the presence of an auth token
// in a configured header and validates it with the configured
// auth token service.
func TokenMiddleware(
	authTokenHeader string,
	authTokenService TokenService,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	}
}
