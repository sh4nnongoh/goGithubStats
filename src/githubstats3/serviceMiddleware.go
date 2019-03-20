package githubstats3

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// ServiceMiddleware describes a service (as opposed to endpoint) middleware.
type ServiceMiddleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingServiceMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) GenerateReport(ctx context.Context, username, token string, repository []string) (repo []Repository, err error) {
	defer func() {
		var str string
		for _, r := range repository {
			str = str + " " + r
		}
		mw.logger.Log("method", "GenerateReport", "username", username, "token", token, "repository", str, "err", err)
	}()
	return mw.next.GenerateReport(ctx, username, token, repository)
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingServiceMiddleware(ints, chars metrics.Counter) ServiceMiddleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  Service
}

func (mw instrumentingMiddleware) GenerateReport(ctx context.Context, username, token string, repository []string) ([]Repository, error) {
	v, err := mw.next.GenerateReport(ctx, username, token, repository)
	mw.ints.Add(float64(len(v)))
	return v, err
}
