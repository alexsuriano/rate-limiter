package repository

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(host, port, pass string, db int) *RedisRepository {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	return &RedisRepository{client: client}
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.client.Get(key).Result()

}

func (r *RedisRepository) Set(key, value string, timeout time.Duration) error {
	return r.client.Set(key, value, timeout).Err()
}

func (r *RedisRepository) Incr(key string) error {
	return r.client.Incr(key).Err()
}
