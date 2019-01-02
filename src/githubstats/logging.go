package githubstats

import (
	"strings"
	"time"

	"github.com/go-kit/kit/log"
)

func NewLoggingMiddleware(logger log.Logger, next GithubReportService) LoggingMiddleware {

	return loggingMiddleware{logger, next}
}

type LoggingMiddleware interface {
	GenerateReport(username, token string, repository []string) (output []repository)
}

type loggingMiddleware struct {
	logger log.Logger
	next   GithubReportService
}

func (mw loggingMiddleware) GenerateReport(username, token string, repository []string) (output []repository) {
	defer func(begin time.Time) {
		var val []string
		for _, r := range output {
			val = append(val, r.RepositoryName)
		}
		mw.logger.Log(
			"method", "generateReport",
			"input", strings.Join(repository[:], ","),
			"output", strings.Join(val[:], ","),
			"err", "nil",
			"took", time.Since(begin),
		)
	}(time.Now())

	output = mw.next.GenerateReport(username, token, repository)
	return
}
