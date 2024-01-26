package middleware

import (
	"net/http"

	"github.com/renanmoreirasan/go-rate-limiter/app/limiter"
	"github.com/renanmoreirasan/go-rate-limiter/infra/configs"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rateLimiter := limiter.RateLimiter{
			Cache: configs.NewCacheRedis(),
		}
		clientInfo, err := rateLimiter.ExtractClientInfoFromRequest(r)
		if err != nil && err.Error() == "IP_NOT_FOUND" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Header X-Real-Ip or X-Forwarded-For not found"))
			return
		}
		if err != nil && err.Error() == "API_KEY_NOT_FOUND" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Header Api-Key not found"))
			return
		}
		if err != nil && err.Error() == "API_KEY_NOT_PERMITTED" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You are not authorized to access this resource. Invalid Api-Key"))
			return
		}
		err = rateLimiter.Check(clientInfo)
		if err != nil && err.Error() == "TOO_MANY_REQUESTS" {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("INTERNAL_SERVER_ERROR"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
