package fetcher

// RepositoryResolver defines a common API for types that allow for resolving of repositories.
type RepositoryResolver interface {
	// ByRepository resolves a repository by its own identifier.
	ByRepository(repository string) (*Repository, error)

	// ByNamespace resolves all repositories for a given namespace identifier.
	ByNamespace(namespace string) (RepositoryList, error)

	// ByUser resolves all repositories for a given user identifier.
	ByUser(user string) (RepositoryList, error)
}

// PipelineResolver defines a common API for types that allow for resolving of pipelines.
type PipelineResolver interface {
	// ByRepository resolves the pipeline for the given repository identifier and branch identifier.
	ByRepository(repository *Repository) (*Pipeline, error)
}
