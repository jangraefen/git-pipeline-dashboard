package fetcher

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	SourceTypeGitlab = "gitlab"
	SourceTypeGithub = "github"
)

var _ RepositoryResolver = &GitlabRepositoryResolver{}

type GitlabRepositoryResolver struct {
	gitlabClient *gitlab.Client
}

func NewGitlabRepositoryResolver(gitlabToken, gitlabEndpoint string) (RepositoryResolver, error) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

	gitlabClient, err := gitlab.NewClient(
		gitlabToken,
		gitlab.WithBaseURL(gitlabEndpoint),
		gitlab.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	return &GitlabRepositoryResolver{gitlabClient: gitlabClient}, nil
}

func (resolver *GitlabRepositoryResolver) ByRepository(repository string) (*Repository, error) {
	project, _, err := resolver.gitlabClient.Projects.GetProject(repository, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	return resolver.toRepository(project), nil
}

func (resolver *GitlabRepositoryResolver) ByNamespace(namespace string) (RepositoryList, error) {
	groupProjects, _, err := resolver.gitlabClient.Groups.ListGroupProjects(namespace, &gitlab.ListGroupProjectsOptions{Simple: gitlab.Bool(true)})
	if err != nil {
		return nil, err
	}

	var repositoryList RepositoryList
	for _, project := range groupProjects {
		repositoryList = append(repositoryList, resolver.toRepository(project))
	}

	return repositoryList, nil
}

func (resolver *GitlabRepositoryResolver) ByUser(user string) (RepositoryList, error) {
	userProjects, _, err := resolver.gitlabClient.Projects.ListUserProjects(user, &gitlab.ListProjectsOptions{Simple: gitlab.Bool(true)})
	if err != nil {
		return nil, err
	}

	var repositoryList RepositoryList
	for _, project := range userProjects {
		repositoryList = append(repositoryList, resolver.toRepository(project))
	}

	return repositoryList, nil
}

func (resolver *GitlabRepositoryResolver) toRepository(project *gitlab.Project) *Repository {
	return &Repository{
		ID:            strconv.Itoa(project.ID),
		Source:        SourceTypeGitlab,
		Name:          project.Name,
		Namespace:     project.Namespace.FullPath,
		DefaultBranch: project.DefaultBranch,
		URL:           project.WebURL,
	}
}

var _ PipelineResolver = &GitlabPipelineResolver{}

type GitlabPipelineResolver struct {
	gitlabClient *gitlab.Client
}

func NewGitlabPipelineResolver(gitlabToken, gitlabEndpoint string) (PipelineResolver, error) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

	gitlabClient, err := gitlab.NewClient(
		gitlabToken,
		gitlab.WithBaseURL(gitlabEndpoint),
		gitlab.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	return &GitlabPipelineResolver{gitlabClient: gitlabClient}, nil
}

func (resolver *GitlabPipelineResolver) ByRepository(repository *Repository) (*Pipeline, error) {
	pipelineInfo, commit, err := getPipelineForRepository(resolver.gitlabClient, repository.ID, repository.DefaultBranch)
	if err != nil {
		return nil, err
	}
	if pipelineInfo == nil {
		return nil, nil
	}

	state := resolver.toPipelineState(pipelineInfo.Status)
	return &Pipeline{
		Ref:           pipelineInfo.Ref,
		URL:           pipelineInfo.WebURL,
		Time:          pipelineInfo.UpdatedAt.Format("02.01.2006 15:04"),
		CommitSHA:     pipelineInfo.SHA,
		CommitAuthor:  commit.AuthorName,
		CommitMessage: strings.Trim(strings.SplitN(commit.Message, "\n", 1)[0], "\n\t "),
		CommitState:   state,
		PipelineRuns: PipelineRunList{{
			Name:  ".gitlab-ci.yml",
			State: state,
			URL:   pipelineInfo.WebURL,
		}},
	}, nil
}

func getPipelineForRepository(git *gitlab.Client, repository, branch string) (*gitlab.PipelineInfo, *gitlab.Commit, error) {
	pipelines, _, err := git.Pipelines.ListProjectPipelines(repository, &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 1},
		Ref:         gitlab.String(branch),
		OrderBy:     gitlab.String("id"),
	})
	if err != nil {
		return nil, nil, err
	}

	if len(pipelines) > 0 {
		commit, _, err := git.Commits.GetCommit(repository, pipelines[0].SHA)
		if err != nil {
			return nil, nil, err
		}

		return pipelines[0], commit, nil
	}

	return nil, nil, nil
}

func (resolver *GitlabPipelineResolver) toPipelineState(status string) PipelineState {
	successStatus := []string{"success"}
	failedStatus := []string{"failed"}
	runningStatus := []string{"created", "waiting_for_resource", "preparing", "pending", "running", "scheduled"}

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

func containsString(slice []string, s string) bool {
	for _, element := range slice {
		if element == s {
			return true
		}
	}

	return false
}
