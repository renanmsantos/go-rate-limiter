package limiter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
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

func (rateLimiter RateLimiter) ExtractClientInfoFromRequest(r *http.Request) ClientInfo {

	key := r.Header.Get("Api-Key")
	if key != "" {
		clientInfo := rateLimiter.getClientInfoFromApiKey(key)
		if clientInfo != (ClientInfo{}) {
			return clientInfo
		}
		key = ""
	}
	if key == "" {
		key = r.Header.Get("X-Real-Ip")
	}
	if key == "" {
		key = r.Header.Get("X-Forwarded-For")
	}

	return ClientInfo{
		Key:             key,
		RequestLimit:    viper.GetInt64("REQUEST_LIMIT"),
		RequestInterval: viper.GetInt64("REQUEST_SECONDS_INTERVAL"),
	}
}

func (rateLimiter RateLimiter) getClientInfoFromApiKey(apiKey string) ClientInfo {
	for _, token := range permittedTokens {
		if apiKey == token.key {
			return ClientInfo{
				Key:             token.key,
				RequestLimit:    token.request_limit,
				RequestInterval: token.request_interval,
			}
		}
	}
	return ClientInfo{}
}

func (rateLimiter RateLimiter) Check(clientInfo ClientInfo) error {

	currentTime := time.Now()
	expires_in := currentTime.Unix() / clientInfo.RequestInterval
	keyWindow := fmt.Sprintf("%s_%d", clientInfo.Key, expires_in)

	rateLimiter.Cache.Incr(clientInfo.Key)
	count, err := rateLimiter.Cache.GetInt64(keyWindow)
	if err != nil && err != redis.Nil {
		log.Printf("Error getting key %s: %s", keyWindow, err.Error())
		return err
	}
	if count >= clientInfo.RequestLimit {
		log.Printf("%s has reached the limit", clientInfo.Key)
		return errors.New("TOO_MANY_REQUESTS")
	}

	rateLimiter.Cache.PipelineIncr(keyWindow)
	expiration := time.Duration(expires_in) * time.Second
	rateLimiter.Cache.PipelineExpire(keyWindow, expiration)
	_, err = rateLimiter.Cache.PipelineExec()
	if err != nil {
		log.Printf(" Error setting key %s: %s", keyWindow, err.Error())
		return err
	}
	log.Printf("IP requested: %s has %d requests.", clientInfo.Key, count+1)
	return nil
}
