package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/jangraefen/git-pipeline-dashboard/fetcher"
	gocache "github.com/patrickmn/go-cache"
)

var repositoryCache = cache.New[*fetcher.Repository](
	store.NewGoCache(gocache.New(
		5*time.Minute,
		10*time.Minute,
	)),
)

func addCachedRepository(repository *fetcher.Repository) error {
	return repositoryCache.Set(context.Background(), getCacheKey(repository.Source, repository.ID), repository)
}

func getCachedRepository(source, id string) (*fetcher.Repository, error) {
	return repositoryCache.Get(context.Background(), getCacheKey(source, id))
}

func getCacheKey(source, id string) string {
	return fmt.Sprintf("%s-%s", source, id)
}
