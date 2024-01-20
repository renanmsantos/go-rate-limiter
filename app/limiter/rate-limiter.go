package limiter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/renanmoreirasan/go-rate-limiter/infra/config"

	"github.com/spf13/viper"
)

type Token struct {
	key              string
	request_limit    int64
	request_interval int64
}

var permittedTokens = []Token{
	{"token-abc", 10, 10},
	{"token-vbb", 5, 10},
	{"token-bvb", 1, 10},
}

func RequestIsPermitted(w http.ResponseWriter, r *http.Request) error {

	cache := config.Cache
	key, request_limit, request_interval := getRequestLimitAndIntervalBasedOnRequest(r)
	currentTime := time.Now()
	expires_in := currentTime.Unix() / request_interval
	keyWindow := fmt.Sprintf("%s_%d", key, expires_in)

	cache.Incr(key)
	count, err := cache.Get(keyWindow).Int64()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting key %s: %s", keyWindow, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if count >= request_limit {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame."))
		log.Printf("%s has reached the limit", key)
		return errors.New("TOO_MANY_REQUESTS")
	}

	pipe := cache.TxPipeline()
	pipe.Incr(keyWindow)
	expiration := time.Duration(expires_in) * time.Second
	pipe.Expire(keyWindow, expiration)
	_, err = pipe.Exec()
	if err != nil {
		log.Printf(" Error setting key %s: %s", keyWindow, err.Error())
		return err
	}
	log.Printf("IP requested: %s has %d requests.", key, count+1)
	return nil
}

func getRequestLimitAndIntervalBasedOnRequest(r *http.Request) (string, int64, int64) {

	requestKey := r.Header.Get("Api-Key")
	if requestKey != "" {
		isValid, limit, interval := apiKeyIsValid(requestKey)
		if isValid {
			return requestKey, limit, interval
		}
		requestKey = ""
	}

	if requestKey == "" {
		requestKey = r.Header.Get("X-Real-Ip")
	}
	if requestKey == "" {
		requestKey = r.Header.Get("X-Forwarded-For")
	}
	if requestKey == "" {
		requestKey = r.RemoteAddr
	}
	return requestKey, viper.GetInt64("REQUEST_LIMIT"), viper.GetInt64("REQUEST_SECONDS_INTERVAL")
}

func apiKeyIsValid(apiKey string) (bool, int64, int64) {
	for _, token := range permittedTokens {
		if apiKey == token.key {
			return true, token.request_limit, token.request_interval
		}
	}
	return false, 0, 0
}
