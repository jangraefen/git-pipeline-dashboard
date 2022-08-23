package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/jangraefen/git-pipeline-dashboard/fetcher"
	"github.com/jangraefen/git-pipeline-dashboard/util"
)

var repositoryCache = util.NewCache[*fetcher.Repository](5 * time.Minute)
var pipelineCache = util.NewCache[*fetcher.Pipeline](1 * time.Minute)

func addCachedRepository(repository *fetcher.Repository) error {
	return repositoryCache.Set(context.Background(), getCacheKey(repository.Source, repository.ID), repository)
}

func getCachedRepository(source, id string) (*fetcher.Repository, error) {
	return repositoryCache.Get(context.Background(), getCacheKey(source, id))
}

func getCacheKey(source, id string) string {
	return fmt.Sprintf("%s-%s", source, id)
}
