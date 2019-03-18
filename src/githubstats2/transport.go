package githubstats2

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func MakeGenerateReportEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(generateReportRequest)
		r := svc.GenerateReport(req.Username, req.Token, req.RepositoryName)
		//if err != nil {
		//	return generateReportRequest{r, err.Error()}, nil
		//}
		return generateReportResponse{r, ""}, nil
	}
}

func DecodeGenerateReportRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request generateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// For each method, we define request and response structs
type generateReportRequest struct {
	Username       string
	Token          string
	RepositoryName []string
}

type generateReportResponse struct {
	Repository []repository
	Err        string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type repository struct {
	RepositoryFullName string
	RepositoryName     string
	CloneURL           string
	LatestCommitDate   string
	LatestCommitAuthor string
}

func (r repository) GetHeaders() []string {
	return []string{"repositoryFullName", "repositoryName", "cloneURL", "latestCommitDate", "latestCommitAuthor"}
}

func (r repository) ToSlice() []string {
	return []string{r.RepositoryFullName, r.RepositoryName, r.CloneURL, r.LatestCommitDate, r.LatestCommitAuthor}
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
