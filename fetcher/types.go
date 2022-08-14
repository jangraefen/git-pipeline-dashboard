package fetcher

import (
	"sort"
	"strings"
)

// PipelineState enumerates all valid states a pipeline can have.
type PipelineState string

const (
	PipelineStateSuccess PipelineState = "success"
	PipelineStateFailed  PipelineState = "failed"
	PipelineStateRunning PipelineState = "running"
	PipelineStateUnknown PipelineState = "unknown"
)

// RepositoryGroup groups several repositories into a named unit.
type RepositoryGroup struct {
	Title        string         `json:"title"`
	Repositories RepositoryList `json:"repositories"`
}

// RepositoryList represents a sortable list of repositories.
type RepositoryList []*Repository

// Repository contains the required metadata of repository that can be queried for pipeline information.
type Repository struct {
	ID            string `json:"id"`
	Source        string `json:"source"`
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	DefaultBranch string `json:"defaultBranch"`
	URL           string `json:"url"`
}

// Pipeline defines all attributes that define a pipeline state.
type Pipeline struct {
	Ref           string        `json:"ref"`
	State         PipelineState `json:"state"`
	URL           string        `json:"url"`
	Time          string        `json:"time"`
	CommitSHA     string        `json:"commitSHA"`
	CommitAuthor  string        `json:"commitAuthor"`
	CommitMessage string        `json:"commitMessage"`
}

var _ sort.Interface = &RepositoryList{}

func (repositoryList RepositoryList) Len() int {
	return len(repositoryList)
}

func (repositoryList RepositoryList) Less(i int, j int) bool {
	elementI, elementJ := repositoryList[i], repositoryList[j]

	namespaceComparison := strings.Compare(elementI.Namespace, elementJ.Namespace)
	if namespaceComparison != 0 {
		return namespaceComparison < 0
	}

	return strings.Compare(elementI.Name, elementJ.Name) < 0
}

func (repositoryList RepositoryList) Swap(i int, j int) {
	repositoryList[i], repositoryList[j] = repositoryList[j], repositoryList[i]
}
