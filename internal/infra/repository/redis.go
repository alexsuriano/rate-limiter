package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(host, port, pass string, db int) (*RedisRepository, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) Get(key string) (int, error) {
	ctx := context.Background()

	value, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return value, nil
}

func (r *RedisRepository) Increment(key string, expire time.Duration) (int, error) {

	ctx := context.Background()

	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if val == 1 {
		r.client.Expire(ctx, key, expire)
	}

	return int(val), nil
}
