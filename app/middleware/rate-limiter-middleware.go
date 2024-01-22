package middleware

import (
	"net/http"

	"github.com/renanmoreirasan/go-rate-limiter/app/limiter"
	"github.com/renanmoreirasan/go-rate-limiter/infra/config"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rateLimiter := limiter.RateLimiter{
			Cache: config.Cache,
		}
		clientInfo := rateLimiter.ExtractClientInfoFromRequest(r)
		if clientInfo.Key == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You are not authorized to access this resource"))
			return
		}
		err := rateLimiter.Check(clientInfo)
		if err != nil && err.Error() == "TOO_MANY_REQUESTS" {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
