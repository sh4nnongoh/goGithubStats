package githubstatscmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
)

func NewGithubReport(username, token string, repositoryList []string) GithubReport {
	var report githubReport
	for r := range GetRepositoryDetails(username, token, repositoryList) {
		report.repo = append(report.repo, r)
	}
	return report
}

type GithubReport interface {
	PrintRepositoryDetails() []string
}

type githubReport struct {
	repo []Repository
}

func (i githubReport) PrintRepositoryDetails() []string {
	w := csv.NewWriter(os.Stdout)
	headers := Repository{}.GetHeaders()
	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	var values []string
	for _, r := range i.repo {
		values = r.ToSlice()
		if err := w.Write(values); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
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

func GetRepositoryDetails(username, token string, repository []string) chan Repository {
	var wg sync.WaitGroup
	c := make(chan Repository)
	wg.Add(len(repository))
	for _, r := range repository {
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
	return c
}
