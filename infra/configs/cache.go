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

func (cache CacheRedis) Incr(key string) {
	cache.client.Incr(key)
}

func (cache CacheRedis) GetInt64(key string) (int64, error) {
	return cache.client.Get(key).Int64()
}

func (cache CacheRedis) TxPipeline() {
	cache.client.TxPipeline()
}

func (cache CacheRedis) PipelineIncr(key string) {
	cache.client.TxPipeline().Incr(key)
}

func (cache CacheRedis) PipelineExpire(keyWindow string, expiration time.Duration) error {
	return cache.client.TxPipeline().Expire(keyWindow, expiration).Err()
}

func (cache CacheRedis) PipelineExec() (interface{}, error) {
	return cache.client.TxPipeline().Exec()
}
