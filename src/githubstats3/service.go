package githubstats3

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

func NewService(logger log.Logger, ints, chars metrics.Counter) Service {
	var svc Service
	{
		//svc = NewService()
		svc = LoggingServiceMiddleware(logger)(service{})
		svc = InstrumentingServiceMiddleware(ints, chars)(svc)
	}
	return svc
}

type Service interface {
	GenerateReport(ctx context.Context, username, token string, repository []string) ([]Repository, error)
}

type service struct{}

func (service) GenerateReport(ctx context.Context, username, token string, repositoryName []string) (repository []Repository, err error) {
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
					//log.Fatalln("error forming HTTP request:", err)
					return err
				}
				req.SetBasicAuth(username, token)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					//log.Fatalln("error obtaining HTTP response:", err)
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
					//log.Fatalln("error forming HTTP request:", err)
					return err
				}
				req.SetBasicAuth(username, token)

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					//log.Fatalln("error obtaining HTTP response:", err)
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
	//var repository []Repository
	for r := range c {
		repository = append(repository, r)
	}
	return repository, nil
}

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)
