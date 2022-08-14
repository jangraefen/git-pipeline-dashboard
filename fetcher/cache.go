package fetcher

import (
	"context"
	"time"

	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	gocache "github.com/patrickmn/go-cache"
)

var repositoryCache = cache.New[*Repository](
	store.NewGoCache(gocache.New(
		5*time.Minute,
		10*time.Minute,
	)),
)

func AddCachedRepository(repository *Repository) error {
	return repositoryCache.Set(context.Background(), repository.ID, repository)
}

func GetCachedRepository(id string) (*Repository, error) {
	return repositoryCache.Get(context.Background(), id)
}
