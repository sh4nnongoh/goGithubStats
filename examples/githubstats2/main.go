package main

import (
	//"log"

	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats2"
	"github.com/sony/gobreaker"
)

func testStateChange(name string, from, to gobreaker.State) {
	fmt.Println(name + " changed ")
}

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()
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

	var svc githubstats2.GithubReportService
	svc = githubstats2.NewGithubReportService()
	svc = githubstats2.NewLoggingMiddleware(logger, svc)
	svc = githubstats2.NewInstrumentingMiddleware(requestCount, requestLatency, countResult, svc)

	// var generateReport endpoint.Endpoint
	// generateReport = githubstats2.MakeGenerateReportEndpoint(svc)
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
	//generateReport = githubstats2.LoggingMiddleware(log.With(logger, "method", "generateReport"))(generateReport)

	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		//httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/generateReport").Handler(httptransport.NewServer(
		githubstats2.MakeGenerateReportEndpoint(svc),
		githubstats2.DecodeGenerateReportRequest,
		githubstats2.EncodeResponse,
		options...,
	))

	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	//http.Handle("/generateReport", generateReportHandler)
	//http.Handle("/metrics", promhttp.Handler())
	//logger.Log("msg", "HTTP", "addr", ":8080")
	//logger.Log("err", http.ListenAndServe(":8080", nil))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: r,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit", <-errs)
}
