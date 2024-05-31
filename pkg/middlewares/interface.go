package middlewares

import "time"

type CacheInterface interface {
	Get(key string) (string, error)
	Set(key, value string, timeout time.Duration) error
	Incr(key string) error
}
