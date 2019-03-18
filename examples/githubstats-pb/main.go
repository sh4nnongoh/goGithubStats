// Example client cmd
// curl -XPOST -d'{"username":"","token":"","repositoryname":["sh4nnongoh/goGithubStats","rubyist/circuitbreaker","twbs/bootstrap","asd/asd12"]}' localhost:8080/generateReport

package main

import (
	//"log"

	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	githubstats "github.com/sh4nnongoh/goGithubStats/src/githubstats-pb"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats-pb/pb"
	"github.com/sony/gobreaker"
)

func testStateChange(name string, from, to gobreaker.State) {
	fmt.Println(name + " changed ")
}

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		httpAddr = fs.String("http.addr", ":8080", "HTTP listen address")
		grpcAddr = fs.String("grpc.addr", ":8002", "Address for gRPC server")
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

	var svc githubstats.Service
	svc = githubstats.NewService()
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
	//generateReport = githubstats2.LoggingMiddleware(log.With(logger, "method", "generateReport"))(generateReport)

	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		//httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/generateReport").Handler(httptransport.NewServer(
		githubstats.MakeGenerateReportEndpoint(svc),
		githubstats.DecodeGenerateReportRequest,
		githubstats.EncodeResponse,
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

	// Transport: gRPC
	go func() {
		//transportLogger := log.NewContext(logger).With("transport", "gRPC")
		level.Info(logger).Log("transport", "GRPC", "addr", *grpcAddr)
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		s := grpc.NewServer() // uses its own, internal context
		pb.RegisterGithubStatsServer(s, githubstats.NewGrpcService(svc))
		_ = logger.Log("addr", *grpcAddr)
		errs <- s.Serve(ln)
	}()

	level.Error(logger).Log("exit", <-errs)
}
