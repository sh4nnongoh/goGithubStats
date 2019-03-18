package githubstats

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

// func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
// 	return func(next endpoint.Endpoint) endpoint.Endpoint {
// 		return func(ctx context.Context, request interface{}) (interface{}, error) {
// 			logger.Log("msg", "calling endpoint")
// 			defer logger.Log("msg", "called endpoint")
// 			return next(ctx, request)
// 		}
// 	}
// }

func NewInstrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	countResult metrics.Histogram,
	next Service) InstrumentingMiddleware {

	return instrumentingMiddleware{requestCount, requestLatency, countResult, next}
}

type InstrumentingMiddleware interface {
	GenerateReport(username, token string, repository []string) (output []repository)
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           Service
}

func (mw instrumentingMiddleware) GenerateReport(username, token string, repository []string) (output []repository) {
	defer func(begin time.Time) {
		lvs := []string{"method", "generateReport", "error", "nil"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output = mw.next.GenerateReport(username, token, repository)
	return
}
