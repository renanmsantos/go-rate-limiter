package configs

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

type CacheRedis struct {
	client *redis.Client
}

func NewCacheRedis() CacheRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})
	return CacheRedis{
		client: client,
	}
}

func (cache CacheRedis) Incr(key string) error {
	return cache.client.Incr(key).Err()
}

func (cache CacheRedis) GetInt64(key string) (int64, error) {
	return cache.client.Get(key).Int64()
}

func (cache CacheRedis) Expire(key string, expiration time.Duration) error {
	return cache.client.Expire(key, expiration).Err()
}
