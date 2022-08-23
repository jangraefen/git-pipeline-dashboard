package util

import (
	"context"
	"time"

	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	gocache "github.com/patrickmn/go-cache"
)

func NewCache[T any](defaultExpiration time.Duration) *cache.Cache[T] {
	return cache.New[T](
		store.NewGoCache(gocache.New(
			defaultExpiration,
			2*defaultExpiration,
		)),
	)
}

func FetchThroughCache[T any, K any](c *cache.Cache[T], key K, fetch func(K) (T, error)) (T, error) {
	value, err := c.Get(context.Background(), key)
	if err == nil {
		return value, nil
	} else if store.NOT_FOUND_ERR != err.Error() {
		// gocritic wants as to do T(nil) instead, which is not possible with generics
		// nolint:gocritic
		return *new(T), err
	}

	value, err = fetch(key)
	if err != nil {
		// gocritic wants as to do T(nil) instead, which is not possible with generics
		// nolint:gocritic
		return *new(T), err
	}

	err = c.Set(context.Background(), key, value)
	if err != nil {
		// gocritic wants as to do T(nil) instead, which is not possible with generics
		// nolint:gocritic
		return *new(T), err
	}

	return value, err
}
