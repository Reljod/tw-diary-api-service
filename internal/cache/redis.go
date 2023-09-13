package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Reljod/tw-diary-api-service/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
	config *config.ConfigSchema
}

type RedisError struct{}

func (e *RedisError) Error() string {
	return "Redis internal error"
}

func (e *RedisError) Code() string {
	return "60001"
}

func CreateRedisCache(config *config.ConfigSchema) *RedisCache {
	db, err := strconv.Atoi(config.Redis.Db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Db name not allowed, %v\n", err)
		return nil
	}

	address := fmt.Sprintf("%s:%v", config.Redis.Host, config.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Redis.Password,
		DB:       db, // use default DB
	})

	cache := &RedisCache{Client: client}
	return cache
}

func (cache *RedisCache) Set(key string, value string, options SetCacheOptions) error {

	var expiry int64
	if options.Expiry != nil {
		expiry = *options.Expiry
	}

	err := cache.Client.Set(context.Background(), key, value, time.Duration(expiry)*time.Second).Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return &RedisError{}
	}

	return nil
}

func (cache *RedisCache) Get(key string) (string, error) {
	val, err := cache.Client.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if err == redis.Nil {
			return "", nil
		}

		return "", &RedisError{}
	}

	return val, nil
}
