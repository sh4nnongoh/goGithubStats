package githubstats

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
)

func NewGithubReportService() GithubReportService {
	return githubReportService{}
}

type GithubReportService interface {
	GenerateReport(username, token string, repository []string) []repository
}

type githubReportService struct{}

func (githubReportService) GenerateReport(username, token string, repositoryName []string) []repository {
	var wg sync.WaitGroup
	c := make(chan repository)
	wg.Add(len(repositoryName))
	for _, r := range repositoryName {
		go func(r string) error {
			defer wg.Done()
			var repo repository
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
	var repository []repository
	for r := range c {
		repository = append(repository, r)
	}
	return repository
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
