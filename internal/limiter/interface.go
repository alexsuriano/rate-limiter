package limiter

import "time"

type LimiterInterface interface {
	OverLimit(register string, reisterType RegisterType) (bool, error)
}

type CacheInterface interface {
	Get(key string) (int, error)
	Increment(key string, expire time.Duration) (int, error)
}
