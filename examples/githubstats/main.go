package main

import (
	//"log"
	"net/http"
	"os"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/sony/gobreaker"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats"
)

func testStateChange(name string, from, to gobreaker.State){
	fmt.Println(name + " changed ")
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	svc := githubstats.NewGithubReportService()

	var generateReport endpoint.Endpoint
	generateReport = githubstats.MakeGenerateReportEndpoint(svc)
	settings := gobreaker.Settings{
		Name: "test",
		MaxRequests: 1,
		Interval: 1000000,
		Timeout: 3000000,
		ReadyToTrip: nil,
		OnStateChange: testStateChange,
	}
	fmt.Println(settings)
	generateReport = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(generateReport)
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
