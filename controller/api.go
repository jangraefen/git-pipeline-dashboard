package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/jangraefen/git-pipeline-dashboard/config"
	"github.com/jangraefen/git-pipeline-dashboard/fetcher"
)

type API struct {
	Selections []config.GroupSelectionConfig

	RepositoryResolvers map[string]fetcher.RepositoryResolver
	PipelineResolvers   map[string]fetcher.PipelineResolver
}

type problem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (controller *API) GetRepositories(writer http.ResponseWriter, request *http.Request) {
	decorateHandler(func() (any, error) {
		var groups []fetcher.RepositoryGroup

		for _, selection := range controller.Selections {
			resolver, ok := controller.RepositoryResolvers[selection.Source]
			if !ok {
				return nil, fmt.Errorf("no resolver registered for source '%s'", selection.Source)
			}

			repositories, err := controller.fetchAllRepositories(resolver, selection)
			if err != nil {
				return nil, err
			}

			for _, repository := range repositories {
				if err := addCachedRepository(repository); err != nil {
					return nil, err
				}
			}

			groups = append(groups, fetcher.RepositoryGroup{
				Title:        selection.Title,
				Repositories: repositories,
			})
		}

		return groups, nil
	})(writer, request)
}

func (controller *API) fetchAllRepositories(resolver fetcher.RepositoryResolver, selection config.GroupSelectionConfig) (fetcher.RepositoryList, error) {
	var repositories fetcher.RepositoryList
	for _, repositoryID := range selection.Repositories {
		repository, err := resolver.ByRepository(repositoryID)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, repository)
	}
	for _, namespace := range selection.Namespaces {
		namespaceRepositories, err := resolver.ByNamespace(namespace)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, namespaceRepositories...)
	}
	for _, user := range selection.Users {
		userRepositories, err := resolver.ByUser(user)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, userRepositories...)
	}

	sort.Sort(repositories)
	return repositories, nil
}

func (controller *API) GetPipelines(writer http.ResponseWriter, request *http.Request) {
	decorateHandler(func() (any, error) {
		repositoryID, ok := mux.Vars(request)["repositoryID"]
		if !ok {
			return nil, fmt.Errorf("missing path parameter repositoryID")
		}

		repositorySource, ok := mux.Vars(request)["repositorySource"]
		if !ok {
			return nil, fmt.Errorf("missing path parameter repositorySource")
		}

		repository, err := getCachedRepository(repositorySource, repositoryID)
		if err != nil {
			return nil, err
		}

		resolver, ok := controller.PipelineResolvers[repository.Source]
		if !ok {
			return nil, fmt.Errorf("no resolver registered for source '%s'", repository.Source)
		}

		return resolver.ByRepository(repository)
	})(writer, request)
}

func decorateHandler(handler func() (any, error)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(true)

		result, err := handler()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(&problem{
				Title:  "Unexpected server error",
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			})

			log.Println("encountered error during handling of request:", err)
			return
		}

		if err := encoder.Encode(result); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(&problem{
				Title:  "Unexpected server error",
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			})

			log.Println("encountered error during handling of request:", err)
			return
		}
	}
}
