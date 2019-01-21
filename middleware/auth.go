package middleware

import (
	"log"
	"net/http"

	"kellnhofer.com/tracker/config"
)

type AuthMiddleware struct {
	conf *config.Config
}

func NewAuthMiddleware(conf *config.Config) *AuthMiddleware {
	return &AuthMiddleware{conf}
}

// --- Public methods ---

func (m AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Authenticated?
	auth := r.Header.Get("Authorization")
	if auth == m.conf.Password {
		// Forward to next handler
		next(w, r)
	} else {
		// Abort
		log.Printf("Unauthorized request! (Authorization header: '%s')", auth)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}
