package githubstats3

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	oldcontext "golang.org/x/net/context"
	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/sh4nnongoh/goGithubStats/src/githubstats3/pb"
)

type grpcServer struct {
	generateReport grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints Set, otTracer stdopentracing.Tracer, logger log.Logger) pb.AddServer {

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		generateReport: grpctransport.NewServer(
			endpoints.GenerateReportEndpoint,
			decodeGRPCGenerateReportRequest,
			encodeGRPCGenerateReportResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "GenerateReport", logger)))...,
		),
	}
}

func (s *grpcServer) GenerateReport(ctx oldcontext.Context, req *pb.GenerateReportRequest) (*pb.GenerateReportResponse, error) {
	_, rep, err := s.GenerateReport.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GenerateReportResponse), nil
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger) Service {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	options := []grpctransport.ClientOption{}

	// Each individual endpoint is an grpc/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var generateReportEndpoint endpoint.Endpoint
	{
		generateReportEndpoint = grpctransport.NewClient(
			conn,
			"pb.GenerateReport",
			"GenerateReport",
			encodeGRPCGenerateReportRequest,
			decodeGRPCGenerateReportResponse,
			pb.GenerateReportReply{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		generateReportEndpoint = opentracing.TraceClient(otTracer, "GenerateReport")(generateReportEndpoint)
		generateReportEndpoint = limiter(generateReportEndpoint)
		generateReportEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GenerateReport",
			Timeout: 30 * time.Second,
		}))(generateReportEndpoint)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return Set{
		GenerateReportEndpoint: generateReportEndpoint,
	}
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGRPCGenerateReportRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GenerateReportRequest)
	return generateReportRequest{Username: req.Username, Token: req.Token, RepositoryName: req.RepositoryName}, nil
}

// decodeGRPCSumResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC sum reply to a user-domain sum response. Primarily useful in a client.
func decodeGRPCGenerateReportResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GenerateReportResponse)
	return generateReportResponse{Repository: reply.Repository, Err: str2err(reply.Err)}, nil
}

// encodeGRPCSumResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain sum response to a gRPC sum reply. Primarily useful in a server.
func encodeGRPCGenerateReportResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(generateReportResponse)
	return &pb.GenerateReportReply{V: int64(resp.V), Err: err2str(resp.Err)}, nil
}

// encodeGRPCSumRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain sum request to a gRPC sum request. Primarily useful in a client.
func encodeGRPCGenerateReportRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(generateReportRequest)
	return &pb.GenerateReportRequest{A: int64(req.A), B: int64(req.B)}, nil
}

// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
