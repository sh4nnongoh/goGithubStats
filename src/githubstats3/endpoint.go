package githubstats3

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	stdopentracing "github.com/opentracing/opentracing-go"
	//stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	GenerateReportEndpoint endpoint.Endpoint
}

// GenerateReport implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) GenerateReport(ctx context.Context, username, token string, repository []string) ([]Repository, error) {
	resp, err := s.GenerateReportEndpoint(ctx, generateReportRequest{Username: username, Token: token, RepositoryName: repository})
	if err != nil {
		return nil, err
	}
	response := resp.(generateReportResponse)
	return response.Repository, response.Err
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func NewSet(svc Service, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer) Set {
	var generateReportEndpoint endpoint.Endpoint
	{
		generateReportEndpoint = MakeGenerateReportEndpoint(svc)
		generateReportEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(generateReportEndpoint)
		generateReportEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(generateReportEndpoint)
		generateReportEndpoint = opentracing.TraceServer(otTracer, "GenerateReport")(generateReportEndpoint)
		generateReportEndpoint = LoggingEndpointMiddleware(log.With(logger, "method", "GenerateReport"))(generateReportEndpoint)
		generateReportEndpoint = InstrumentingEndpointMiddleware(duration.With("method", "GenerateReport"))(generateReportEndpoint)
	}
	return Set{
		GenerateReportEndpoint: generateReportEndpoint,
	}
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeGenerateReportEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(generateReportRequest)
		v, err := s.GenerateReport(ctx, req.Username, req.Token, req.RepositoryName)
		return generateReportResponse{Repository: v, Err: err}, nil
	}
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

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = generateReportResponse{}
)

// For each method, we define request and response structs
type generateReportRequest struct {
	Username       string
	Token          string
	RepositoryName []string
}

type generateReportResponse struct {
	Repository []Repository
	Err        error `json:"err,omitempty"` // errors don't define JSON marshaling
}

// Failed implements endpoint.Failer.
func (r generateReportResponse) Failed() error { return r.Err }

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
