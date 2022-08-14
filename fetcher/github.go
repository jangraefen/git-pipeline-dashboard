package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

var _ RepositoryResolver = &GithubRepositoryResolver{}

type GithubRepositoryResolver struct {
	githubClient *github.Client
}

func NewGithubRepositoryResolver(githubToken string) (RepositoryResolver, error) {
	httpClient := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	ctx := context.WithValue(context.TODO(), oauth2.HTTPClient, httpClient)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tokenHTTPClient := oauth2.NewClient(ctx, tokenSource)
	tokenHTTPClient.Timeout = 5 * time.Second

	return &GithubRepositoryResolver{githubClient: github.NewClient(tokenHTTPClient)}, nil
}

func (resolver *GithubRepositoryResolver) ByNamespace(namespace string) (RepositoryList, error) {
	repositories, _, err := resolver.githubClient.Repositories.ListByOrg(context.TODO(), namespace, nil)
	if err != nil {
		return nil, err
	}

	var repositoryList RepositoryList
	for _, repository := range repositories {
		repositoryList = append(repositoryList, resolver.toRepository(repository))
	}

	return repositoryList, nil
}

func (resolver *GithubRepositoryResolver) ByRepository(repository string) (*Repository, error) {
	if strings.ContainsRune(repository, '/') {
		parts := strings.SplitN(repository, "/", 2)
		repo, _, err := resolver.githubClient.Repositories.Get(context.TODO(), parts[0], parts[1])
		if err != nil {
			return nil, err
		}

		return resolver.toRepository(repo), nil
	} else if repositoryID, err := strconv.ParseInt(repository, 10, 0); err == nil {
		repo, _, err := resolver.githubClient.Repositories.GetByID(context.TODO(), repositoryID)
		if err != nil {
			return nil, err
		}

		return resolver.toRepository(repo), nil
	}

	return nil, fmt.Errorf("could not identify format of repository ID '%s'", repository)
}

func (resolver *GithubRepositoryResolver) ByUser(user string) (RepositoryList, error) {
	repositories, _, err := resolver.githubClient.Repositories.List(context.TODO(), user, nil)
	if err != nil {
		return nil, err
	}

	var repositoryList RepositoryList
	for _, repository := range repositories {
		repositoryList = append(repositoryList, resolver.toRepository(repository))
	}

	return repositoryList, nil
}

func (resolver *GithubRepositoryResolver) toRepository(repository *github.Repository) *Repository {
	return &Repository{
		ID:            strconv.FormatInt(repository.GetID(), 10),
		Source:        SourceTypeGithub,
		Name:          repository.GetName(),
		Namespace:     repository.GetOwner().GetLogin(),
		DefaultBranch: repository.GetDefaultBranch(),
		URL:           repository.GetHTMLURL(),
	}
}

var _ PipelineResolver = &GithubPipelineResolver{}

type GithubPipelineResolver struct {
	githubClient *github.Client
}

func NewGithubPipelineResolver(githubToken string) (PipelineResolver, error) {
	httpClient := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	ctx := context.WithValue(context.TODO(), oauth2.HTTPClient, httpClient)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tokenHTTPClient := oauth2.NewClient(ctx, tokenSource)
	tokenHTTPClient.Timeout = 5 * time.Second

	return &GithubPipelineResolver{githubClient: github.NewClient(tokenHTTPClient)}, nil
}

func (resolver *GithubPipelineResolver) ByRepository(repository *Repository) (*Pipeline, error) {
	workflowRuns, _, err := resolver.githubClient.Actions.ListRepositoryWorkflowRuns(context.TODO(), repository.Namespace, repository.Name, &github.ListWorkflowRunsOptions{Branch: repository.DefaultBranch})
	if err != nil {
		return nil, err
	} else if workflowRuns.GetTotalCount() == 0 {
		return nil, nil
	}

	// TODO: This is still very much flawed, since we do not know, which workflow is actually the one we are interessted in?
	var latestWorkflowRun *github.WorkflowRun
	for _, workflowRun := range workflowRuns.WorkflowRuns {
		if latestWorkflowRun == nil || workflowRun.GetCreatedAt().After(latestWorkflowRun.GetCreatedAt().Time) {
			latestWorkflowRun = workflowRun
		}
	}

	return &Pipeline{
		Ref:           repository.DefaultBranch,
		State:         resolver.toPipelineState(latestWorkflowRun.GetStatus()),
		URL:           latestWorkflowRun.GetHTMLURL(),
		Time:          latestWorkflowRun.GetUpdatedAt().Format("02.01.2006 15:04"),
		CommitSHA:     latestWorkflowRun.GetHeadCommit().GetSHA(),
		CommitAuthor:  latestWorkflowRun.GetHeadCommit().GetAuthor().GetName(),
		CommitMessage: strings.Trim(strings.SplitN(latestWorkflowRun.GetHeadCommit().GetMessage(), "\n", 1)[0], "\n\t "),
	}, nil
}

func (resolver *GithubPipelineResolver) toPipelineState(status string) PipelineState {
	successStatus := []string{"completed", "success"}
	failedStatus := []string{"failure", "timed_out"}
	runningStatus := []string{"action_required", "in_progress", "queued", "requested", "waiting"}

	switch {
	case containsString(successStatus, status):
		return PipelineStateSuccess
	case containsString(failedStatus, status):
		return PipelineStateFailed
	case containsString(runningStatus, status):
		return PipelineStateRunning
	default:
		return PipelineStateUnknown
	}
}
