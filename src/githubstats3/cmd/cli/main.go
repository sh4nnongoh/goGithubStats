package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"google.golang.org/grpc"

	lightstep "github.com/lightstep/lightstep-tracer-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"

	"github.com/go-kit/kit/log"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats3"
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an githubstats2. This presumption is reflected in
	// the addcli binary and the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.
	fs := flag.NewFlagSet("addcli", flag.ExitOnError)
	var (
		httpAddr       = fs.String("http-addr", "", "HTTP address of githubstats2")
		grpcAddr       = fs.String("grpc-addr", "", "gRPC address of githubstats2")
		jsonRPCAddr    = fs.String("jsonrpc-addr", "", "JSON RPC address of githubstats2")
		lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
		appdashAddr    = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
		method         = fs.String("method", "generateReport", "generateReport")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <username> <token> <repositoryName>")
	fs.Parse(os.Args[1:])
	if len(fs.Args()) <= 2 {
		fmt.Println("Arg error exiting...")
		fs.Usage()
		os.Exit(1)
	}

	// This is a demonstration client, which supports multiple tracers.
	// Your clients will probably just use one tracer.
	var otTracer stdopentracing.Tracer
	{
		if *lightstepToken != "" {
			otTracer = lightstep.NewTracer(lightstep.Options{
				AccessToken: *lightstepToken,
			})
			defer lightstep.FlushLightStepTracer(otTracer)
		} else if *appdashAddr != "" {
			otTracer = appdashot.NewTracer(appdash.NewRemoteCollector(*appdashAddr))
		} else {
			otTracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.
	var (
		svc githubstats3.Service
		err error
	)
	if *httpAddr != "" {
		svc, err = githubstats3.NewHTTPClient(*httpAddr, otTracer, log.NewNopLogger())
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		//svc = githubstats3.NewGRPCClient(conn, otTracer, log.NewNopLogger())
	} else if *jsonRPCAddr != "" {
		//svc, err = githubstats3.NewJSONRPCClient(*jsonRPCAddr, otTracer, log.NewNopLogger())
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch *method {
	case "generateReport":
		username := fs.Args()[0]
		token := fs.Args()[1]
		repository := fs.Args()[2:]
		v, err := svc.GenerateReport(context.Background(), username, token, repository)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%d", len(v))

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
