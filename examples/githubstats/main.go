package main

import (
	//"log"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	svc := githubstats.NewGithubReportService()

	var generateReport endpoint.Endpoint
	generateReport = githubstats.MakeGenerateReportEndpoint(svc)
	generateReport = githubstats.LoggingMiddleware(log.With(logger, "method", "generateReport"))(generateReport)

	generateReportHandler := httptransport.NewServer(
		generateReport,
		githubstats.DecodeGenerateReportRequest,
		githubstats.EncodeResponse,
	)

	http.Handle("/generateReport", generateReportHandler)
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
