package githubstatswithextpkg

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/sh4nnongoh/goGithubStats/src/githubstatsterminal"
)

func NewService() Service {
	return Service{}
}

type Service interface {
	GenerateReport(username, token string, repository []string) string
}

type Service struct{}

func (Service) GenerateReport(username, token string, repositoryName []string) string {
	report := githubstatsterminal.NewGithubReport(username, token, repositoryName)
	return report.PrintRepositoryDetailsString()
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
	Report string
	Err    string `json:"err,omitempty"` // errors don't define JSON marshaling
}

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
