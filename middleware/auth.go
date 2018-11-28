package middleware

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"

	"kellnhofer.com/tracker/config"
)

type AuthMiddleware struct {
	conf *config.Config
}

func NewAuthMiddleware(conf *config.Config) *AuthMiddleware {
	return &AuthMiddleware{conf}
}

// --- Public methods ---

func (m AuthMiddleware) GetAuthHandlerFunc() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		m.checkAuth(w, r, next)
	}
}

// --- Private methods ---

func (m AuthMiddleware) checkAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
