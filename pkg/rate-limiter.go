package pkg

import (
	"log"
	"net/http"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
