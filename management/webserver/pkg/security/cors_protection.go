package security

import (
	"net/http"
)

// CORSEnforcer is a middleware to enforce CORS security policies
func CORSEnforcer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate Origin/Referer
		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")

		if !isValidOrigin(origin) && !isValidReferer(referer) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isValidOrigin checks the provided origin
func isValidOrigin(origin string) bool {
	// Implement your validation logic
	return origin == "https://example.com"
}

// isValidReferer checks the provided referer
func isValidReferer(referer string) bool {
	// Implement your validation logic
	return referer == "https://example.com/"
}