package limiter

import "time"

type CacheInterface interface {
	Incr(key string)
	GetInt64(key string) (int64, error)
	TxPipeline()
	PipelineIncr(key string)
	PipelineExpire(keyWindow string, expiration time.Duration) error
	PipelineExec() (interface{}, error)
}
