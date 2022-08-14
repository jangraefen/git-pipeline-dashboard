package fetcher

import (
	"sort"
	"strings"
)

const (
	PipelineStateSuccess PipelineState = "success"
	PipelineStateFailed  PipelineState = "failed"
	PipelineStateRunning PipelineState = "running"
	PipelineStateUnknown PipelineState = "unknown"
)

// PipelineState enumerates all valid states a pipeline can have.
type PipelineState string

func (state PipelineState) ToInt() int {
	switch state {
	case PipelineStateFailed:
		return 3
	case PipelineStateRunning:
		return 2
	case PipelineStateSuccess:
		return 1
	default:
		return 0
	}
}

func (state PipelineState) CompareTo(other PipelineState) int {
	return state.ToInt() - other.ToInt()
}

// RepositoryGroup groups several repositories into a named unit.
type RepositoryGroup struct {
	Title        string         `json:"title"`
	Repositories RepositoryList `json:"repositories"`
}

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
	Ref           string          `json:"ref"`
	URL           string          `json:"url"`
	Time          string          `json:"time"`
	CommitSHA     string          `json:"commitSHA"`
	CommitAuthor  string          `json:"commitAuthor"`
	CommitMessage string          `json:"commitMessage"`
	CommitState   PipelineState   `json:"commitState"`
	PipelineRuns  PipelineRunList `json:"pipelineRuns"`
}

type PipelineRun struct {
	Name  string        `json:"name"`
	State PipelineState `json:"state"`
	URL   string        `json:"url"`
}

var _ sort.Interface = &RepositoryList{}

// RepositoryList represents a sortable list of repositories.
type RepositoryList []*Repository

func (list RepositoryList) Len() int {
	return len(list)
}

func (list RepositoryList) Less(i int, j int) bool {
	elementI, elementJ := list[i], list[j]

	namespaceComparison := strings.Compare(elementI.Namespace, elementJ.Namespace)
	if namespaceComparison != 0 {
		return namespaceComparison < 0
	}

	return strings.Compare(elementI.Name, elementJ.Name) < 0
}

func (list RepositoryList) Swap(i int, j int) {
	list[i], list[j] = list[j], list[i]
}

var _ sort.Interface = &PipelineRunList{}

// PipelineRunList represents a sortable list of pipeline runs, that all belong to the same commit.
type PipelineRunList []PipelineRun

func (list PipelineRunList) GetState() PipelineState {
	state := PipelineStateUnknown
	for _, run := range list {
		if run.State.CompareTo(state) > 0 {
			state = run.State
		}
	}

	return state
}

func (list PipelineRunList) Len() int {
	return len(list)
}

func (list PipelineRunList) Less(i int, j int) bool {
	return strings.Compare(list[i].Name, list[j].Name) < 0
}

func (list PipelineRunList) Swap(i int, j int) {
	list[i], list[j] = list[j], list[i]
}
