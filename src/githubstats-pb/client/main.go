package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	grpcclient "github.com/sh4nnongoh/goGithubstats/src/githubstats-pb/client/grpc"
	//httpjsonclient "github.com/sh4nnongoh/goGithubstats/src/githubstats-pb/client/httpjson"
	"github.com/go-kit/kit/log"
	githubstats "github.com/sh4nnongoh/goGithubstats/src/githubstats-pb"
)

func main() {
	var (
		transport = flag.String("transport", "httpjson", "httpjson, grpc, netrpc, thrift")
		httpAddr  = flag.String("http.addr", "localhost:8001", "Address for HTTP (JSON) server")
		grpcAddr  = flag.String("grpc.addr", "localhost:8002", "Address for gRPC server")
	)
	flag.Parse()
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "\n%s [flags] method arg1 arg2 arg3\n\n", filepath.Base(os.Args[0]))
		flag.Usage()
		os.Exit(1)
	}

	root := context.Background()
	method, username, token, repository := flag.Arg(0), flag.Arg(1), flag.Arg(2), flag.Arg(3)

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	logger = log.NewContext(logger).With("transport", *transport)

	var svc githubstats.Service
	switch *transport {
	case "grpc":
		cc, err := grpc.Dial(*grpcAddr)
		if err != nil {
			_ = logger.Log("err", err)
			os.Exit(1)
		}
		defer cc.Close()
		svc = grpcclient.New(root, cc, logger)

	case "httpjson": /*
			rawurl := *httpAddr
			if !strings.HasPrefix("http", rawurl) {
				rawurl = "http://" + rawurl
			}
			baseurl, err := url.Parse(rawurl)
			if err != nil {
				_ = logger.Log("err", err)
				os.Exit(1)
			}
			svc = httpjsonclient.New(root, baseurl, logger, nil)
		*/
	default:
		_ = logger.Log("err", "invalid transport")
		os.Exit(1)
	}

	begin := time.Now()
	switch method {
	case "generateReport":
		response := svc.GenerateReport(username, token, repository)
		_ = logger.Log("method", "generateReport", "username", username, "token", token, "", repository, "response", response, "took", time.Since(begin))

	default:
		_ = logger.Log("err", "invalid method "+method)
		os.Exit(1)
	}
}
