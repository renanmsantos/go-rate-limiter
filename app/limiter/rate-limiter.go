package limiter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type RateLimiter struct {
	Cache CacheInterface
}

type ClientInfo struct {
	Key             string
	RequestLimit    int64
	RequestInterval int64
}

func (rateLimiter RateLimiter) ExtractClientInfoFromRequest(r *http.Request) (ClientInfo, error) {
	header := r.Header
	if viper.GetString("REQUEST_LIMITER_MODE") == "IP" {
		key := header.Get("X-Real-Ip")
		if key == "" {
			key = r.Header.Get("X-Forwarded-For")
		}
		if key == "" {
			key = r.RemoteAddr
		}
		return ClientInfo{
			Key:             key,
			RequestLimit:    viper.GetInt64("REQUEST_LIMIT"),
			RequestInterval: viper.GetInt64("REQUEST_SECONDS_INTERVAL"),
		}, nil
	} else if viper.GetString("REQUEST_LIMITER_MODE") == "API_KEY" || viper.GetString("REQUEST_LIMITER_MODE") == "" {
		key := header.Get("Api-Key")
		if key == "" {
			return ClientInfo{}, errors.New("API_KEY_NOT_FOUND")
		}
		clientInfo, err := rateLimiter.validateIfApiKeyIsPermitted(key)
		return clientInfo, err
	}
	return ClientInfo{}, errors.New("REQUEST_LIMITER_MODE_NOT_FOUND")
}

func (rateLimiter RateLimiter) validateIfApiKeyIsPermitted(apiKey string) (ClientInfo, error) {
	for _, token := range permittedTokens {
		if apiKey == token.key {
			return ClientInfo{
				Key:             token.key,
				RequestLimit:    token.request_limit,
				RequestInterval: token.request_interval,
			}, nil
		}
	}
	return ClientInfo{}, errors.New("API_KEY_NOT_PERMITTED")
}

func (rateLimiter RateLimiter) Check(clientInfo ClientInfo) error {

	currentTime := time.Now()
	expires_in := currentTime.Unix() / clientInfo.RequestInterval
	keyWindow := fmt.Sprintf("%s_%d", clientInfo.Key, expires_in)

	err := rateLimiter.Cache.Incr(clientInfo.Key)
	if err != nil {
		return err
	}
	count, _ := rateLimiter.Cache.GetInt64(keyWindow)
	if count >= clientInfo.RequestLimit {
		log.Printf("%s has reached the limit", clientInfo.Key)
		return errors.New("TOO_MANY_REQUESTS")
	}

	rateLimiter.Cache.Incr(keyWindow)
	expiration := time.Duration(expires_in) * time.Second
	err = rateLimiter.Cache.Expire(keyWindow, expiration)
	if err != nil {
		log.Printf(" Error setting key %s: %s", keyWindow, err.Error())
		return err
	}
	log.Printf("IP requested: %s has %d requests.", clientInfo.Key, count+1)
	return nil
}
