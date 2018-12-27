package githubstats

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/go-kit/kit/endpoint"
)

func NewGithubReportService() GithubReportService {
	return githubReportService{}
}

type GithubReportService interface {
	GenerateReport(username, token string, repository []string) []Repository
}

type githubReportService struct{}

func (githubReportService) GenerateReport(username, token string, repositoryName []string) []Repository {
	var wg sync.WaitGroup
	c := make(chan Repository)
	wg.Add(len(repositoryName))
	for _, r := range repositoryName {
		go func(r string) error {
			defer wg.Done()
			var repo Repository
			repo.RepositoryFullName = r
			{
				req, err := http.NewRequest("GET", "https://api.github.com/repos/"+r, nil)
				if err != nil {
					log.Fatalln("error forming HTTP request:", err)
					return err
				}
				req.SetBasicAuth(username, token)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalln("error obtaining HTTP response:", err)
					return err
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					return errors.New(resp.Status)
				}

				var jRsp repositoryResponse
				json.NewDecoder(resp.Body).Decode(&jRsp)

				repo.RepositoryName = jRsp.Name
				repo.CloneURL = jRsp.Clone_url
			}
			{
				req, err := http.NewRequest("GET", "https://api.github.com/repos/"+r+"/commits", nil)
				if err != nil {
					log.Fatalln("error forming HTTP request:", err)
					return err
				}
				req.SetBasicAuth(username, token)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatalln("error obtaining HTTP response:", err)
					return err
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					return errors.New(resp.Status)
				}

				var jRsp []repositoryCommitsResponse
				json.NewDecoder(resp.Body).Decode(&jRsp)

				repo.LatestCommitDate = jRsp[0].Commit.Committer.Date
				repo.LatestCommitAuthor = jRsp[0].Commit.Author.Name
			}
			c <- repo
			return nil
		}(r)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	var repository []Repository
	for r := range c {
		repository = append(repository, r)
	}
	return repository
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// For each method, we define request and response structs
type generateReportRequest struct {
	Username       string
	Token          string
	RepositoryName []string
}

type generateReportResponse struct {
	Repository []Repository
	Err        string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func MakeGenerateReportEndpoint(svc GithubReportService) endpoint.Endpoint {
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

type Repository struct {
	RepositoryFullName string
	RepositoryName     string
	CloneURL           string
	LatestCommitDate   string
	LatestCommitAuthor string
}

func (r Repository) GetHeaders() []string {
	return []string{"repositoryFullName", "repositoryName", "cloneURL", "latestCommitDate", "latestCommitAuthor"}
}

func (r Repository) ToSlice() []string {
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
