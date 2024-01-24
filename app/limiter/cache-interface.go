package limiter

import "time"

type CacheInterface interface {
	Incr(key string) error
	GetInt64(key string) (int64, error)
	Expire(key string, expiration time.Duration) error
}
