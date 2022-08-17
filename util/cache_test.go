package util

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/eko/gocache/v3/store"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestNewCache(t *testing.T) {
	sutCache := NewCache[string](50 * time.Millisecond)

	err := sutCache.Set(context.Background(), "key", "value")
	assert.NoError(t, err)

	val, err := sutCache.Get(context.Background(), "key")
	assert.Equal(t, "value", val)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	_, err = sutCache.Get(context.Background(), "key")
	assert.Error(t, err)
	assert.Equal(t, store.NOT_FOUND_ERR, err.Error())
}

func TestFetchThroughCache(t *testing.T) {
	testCache := NewCache[string](50 * time.Millisecond)

	testFetcherErrorToggle := atomic.NewBool(true)
	testFetcher := func(key string) (string, error) {
		if testFetcherErrorToggle.Load() {
			return "", fmt.Errorf("could not fetch value")
		}

		return key, nil
	}

	testCache.Set(context.Background(), "value1", "value1")
	val, err := FetchThroughCache(testCache, "value1", testFetcher)
	assert.Equal(t, "value1", val)
	assert.NoError(t, err)

	testFetcherErrorToggle.Store(false)
	val, err = FetchThroughCache(testCache, "value2", testFetcher)
	assert.Equal(t, "value2", val)
	assert.NoError(t, err)

	val, err = testCache.Get(context.Background(), "value2")
	assert.Equal(t, "value2", val)
	assert.NoError(t, err)

	testFetcherErrorToggle.Store(true)
	_, err = FetchThroughCache(testCache, "err", testFetcher)
	assert.Error(t, err)
	testFetcherErrorToggle.Store(false)
}
