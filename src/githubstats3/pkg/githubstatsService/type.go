package githubstatsService

type Repository struct {
	RepositoryFullName string
	RepositoryName     string
	CloneURL           string
	LatestCommitDate   string
	LatestCommitAuthor string
}

// curl -i https://api.github.com/repos/<org>/<name>
type repositoryResponse struct {
	Id          int
	Node_id     string
	Name        string
	Full_name   string
	Description string
	Clone_url   string
}

// curl -i https://api.github.com/repos/<org>/<name>/commits
type repositoryCommitsResponse struct {
	Commit struct {
		Author struct {
			Name  string
			Email string
			Date  string
		}
		Committer struct {
			Name  string
			Email string
			Date  string
		}
		Message string
	}
	Url string
}
