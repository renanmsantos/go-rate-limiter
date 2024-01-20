package middleware

import (
	"net/http"

	"github.com/renanmoreirasan/go-rate-limiter/app/limiter"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := limiter.RequestIsPermitted(w, r)
		if err != nil {
			return
		}
		next.ServeHTTP(w, r)
	})
}
