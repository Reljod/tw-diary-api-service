package cache

import (
	"fmt"
	"os"
)

type SessionCache interface {
	Set(key string, value string, overrideOpts *SessionCacheOptions) error
	Get(key string) (string, error)
	Delete(key string) error
}

type SessionCacheOptions struct {
	Expiry int64
	Prefix string
}

type SessionRedisCache struct {
	Redis   *RedisCache
	Options *SessionCacheOptions
}

func (cache *SessionRedisCache) Set(key string, value string, overrideOpts *SessionCacheOptions) error {
	var expiry int64 = cache.Options.Expiry
	var prefix string = cache.Options.Prefix

	if overrideOpts != nil {
		expiry = overrideOpts.Expiry
		prefix = overrideOpts.Prefix
	}

	keyWithPrefix := cache.buildKey(key, &prefix)

	options := SetCacheOptions{Expiry: &expiry}
	err := cache.Redis.Set(keyWithPrefix, value, options)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting session cache: %v\n", err)
		return err
	}

	return nil
}

func (cache *SessionRedisCache) Get(key string) (string, error) {
	keyWithPrefix := cache.buildKey(key, nil)

	val, err := cache.Redis.Get(keyWithPrefix)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting session cache: %v\n", err)
		return "", err
	}

	return val, nil
}

func (cache *SessionRedisCache) Delete(key string) error {
	keyWithPrefix := cache.buildKey(key, nil)
	err := cache.Redis.Delete(keyWithPrefix)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting session cache: %v\n", err)
		return err
	}

	return nil
}

func (cache *SessionRedisCache) buildKey(key string, prefix *string) string {
	if prefix != nil {
		return fmt.Sprintf("%s%s", *prefix, key)
	}
	return fmt.Sprintf("%s%s", cache.Options.Prefix, key)
}
