package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jangraefen/git-pipeline-dashboard/config"
	"github.com/jangraefen/git-pipeline-dashboard/controller"
	"github.com/jangraefen/git-pipeline-dashboard/fetcher"
	"github.com/jangraefen/git-pipeline-dashboard/util"
)

var (
	serverAddr     string
	serverPort     uint
	gitlabEndpoint string
	gitlabToken    string
	githubToken    string

	//go:embed frontend/dist
	frontend embed.FS
)

func main() {
	serverAddrDefault := util.GetEnvOrDefault(fmt.Sprintf("%s_SERVER_ADDRESS", config.EnvironmentVariablePrefix), "0.0.0.0")
	flag.StringVar(&serverAddr, "server-address", serverAddrDefault, "The address the server will be bound to")
	serverPortDefault := util.GetEnvOrDefaultUint(fmt.Sprintf("%s_SERVER_PORT", config.EnvironmentVariablePrefix), 8080)
	flag.UintVar(&serverPort, "server-port", serverPortDefault, "The port the server will be bound to")
	gitlabEndpointDefault := util.GetEnvOrDefault(fmt.Sprintf("%s_GITLAB_ENDPOINT", config.EnvironmentVariablePrefix), "https://gitlab.com/api/v4")
	flag.StringVar(&gitlabEndpoint, "gitlab-endpoint", gitlabEndpointDefault, "The GitLab API endpoint to read from")
	gitlabTokenDefault := util.GetEnvOrDefault(fmt.Sprintf("%s_GITLAB_TOKEN", config.EnvironmentVariablePrefix), "")
	flag.StringVar(&gitlabToken, "gitlab-token", gitlabTokenDefault, "The GitLab access token to authenticate to the API endpoint")
	githubTokenDefault := util.GetEnvOrDefault(fmt.Sprintf("%s_GITHUB_TOKEN", config.EnvironmentVariablePrefix), "")
	flag.StringVar(&githubToken, "github-token", githubTokenDefault, "The Github access token to authenticate to the API endpoint")
	flag.Parse()

	selections := config.FromEnvironmentVariables()
	repositoryResolvers := make(map[string]fetcher.RepositoryResolver, 0)
	pipelineResolvers := make(map[string]fetcher.PipelineResolver, 0)

	if selections.ContainsSource(fetcher.SourceTypeGitlab) {
		gitlabRepositoryResolver, gitlabPipelineResolver, err := getGitlabSupport(gitlabToken, gitlabEndpoint)
		if err != nil {
			log.Fatalln("unable to create gitlab client:", err)
			return
		}
		repositoryResolvers[fetcher.SourceTypeGitlab] = gitlabRepositoryResolver
		pipelineResolvers[fetcher.SourceTypeGitlab] = gitlabPipelineResolver
	}

	if selections.ContainsSource(fetcher.SourceTypeGithub) {
		githubRepositoryResolver, githubPipelineResolver, err := getGithubSupport(githubToken)
		if err != nil {
			log.Fatalln("unable to create github client:", err)
			return
		}
		repositoryResolvers[fetcher.SourceTypeGithub] = githubRepositoryResolver
		pipelineResolvers[fetcher.SourceTypeGithub] = githubPipelineResolver
	}

	apiController := &controller.API{
		Selections:          selections,
		RepositoryResolvers: repositoryResolvers,
		PipelineResolvers:   pipelineResolvers,
	}

	sub, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatalln("unable to start http server:", err)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/repositories", apiController.GetRepositories)
	router.HandleFunc("/repositories/{repositorySource}/{repositoryID}", apiController.GetPipelines)
	router.PathPrefix("/").Handler(http.FileServer(http.FS(sub)))

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", serverAddr, serverPort), router); err != nil {
		log.Fatalln("unable to start http server:", err)
		return
	}
}

func getGitlabSupport(token, endpoint string) (fetcher.RepositoryResolver, fetcher.PipelineResolver, error) {
	if token == "" {
		return nil, nil, fmt.Errorf("missing '%s_GITLAB_TOKEN' environment variable or 'gitlab-token' parameter", config.EnvironmentVariablePrefix)
	}

	repositoryResolver, err := fetcher.NewGitlabRepositoryResolver(token, endpoint)
	if err != nil {
		return nil, nil, err
	}
	pipelineResolver, err := fetcher.NewGitlabPipelineResolver(token, endpoint)
	if err != nil {
		return nil, nil, err
	}

	return repositoryResolver, pipelineResolver, nil
}

func getGithubSupport(token string) (fetcher.RepositoryResolver, fetcher.PipelineResolver, error) {
	if token == "" {
		return nil, nil, fmt.Errorf("missing '%s_GITHUB_TOKEN' environment variable or 'github-token' parameter", config.EnvironmentVariablePrefix)
	}

	repositoryResolver, err := fetcher.NewGithubRepositoryResolver(token)
	if err != nil {
		return nil, nil, err
	}
	pipelineResolver, err := fetcher.NewGithubPipelineResolver(token)
	if err != nil {
		return nil, nil, err
	}

	return repositoryResolver, pipelineResolver, nil
}
