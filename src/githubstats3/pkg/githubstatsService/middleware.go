package githubstatsService

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
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
		mw.logger.Log("method", "GenerateReport", "username", username, "token", token, "repository", repository, "err", err)
	}()
	return mw.next.GenerateReport(ctx, username, token, repository)
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
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
