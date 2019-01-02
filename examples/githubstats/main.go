package main

import (
	//"log"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats"
	"github.com/sony/gobreaker"
)

func testStateChange(name string, from, to gobreaker.State) {
	fmt.Println(name + " changed ")
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "github_report_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "github_report_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "github_report_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	var svc githubstats.GithubReportService
	svc = githubstats.NewGithubReportService()
	svc = githubstats.NewLoggingMiddleware(logger, svc)
	svc = githubstats.NewInstrumentingMiddleware(requestCount, requestLatency, countResult, svc)

	// var generateReport endpoint.Endpoint
	// generateReport = githubstats.MakeGenerateReportEndpoint(svc)
	// settings := gobreaker.Settings{
	// 	Name:          "test",
	// 	MaxRequests:   1,
	// 	Interval:      1000000,
	// 	Timeout:       3000000,
	// 	ReadyToTrip:   nil,
	// 	OnStateChange: testStateChange,
	// }
	// fmt.Println(settings)
	//generateReport = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(generateReport)
	//generateReport = githubstats.LoggingMiddleware(log.With(logger, "method", "generateReport"))(generateReport)

	generateReportHandler := httptransport.NewServer(
		githubstats.MakeGenerateReportEndpoint(svc),
		githubstats.DecodeGenerateReportRequest,
		githubstats.EncodeResponse,
	)

	http.Handle("/generateReport", generateReportHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
