package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	githubstats "github.com/sh4nnongoh/goGithubStats/src/githubstats-pb"
	"github.com/sh4nnongoh/goGithubStats/src/githubstats-pb/pb"
)

// New returns a GithubStatsService that's backed by the provided ClientConn.
func New(ctx context.Context, cc *grpc.ClientConn, logger log.Logger) githubstats.Service {
	return client{ctx, pb.NewGithubStatsClient(cc), logger}
}

type client struct {
	context.Context
	pb.GithubStatsClient
	log.Logger
}

// TODO(pb): If your service interface methods don't return an error, we have
// no way to signal problems with a service client. If they don't take a
// context, we have to provide a global context for any transport that
// requires one, effectively making your service a black box to any context-
// specific information. So, we should make some recommendations:
//
// - To get started, a simple service interface is probably fine.
//
// - To properly deal with transport errors, every method on your service
//   should return an error. This is probably important.
//
// - To properly deal with context information, every method on your service
//   can take a context as its first argument. This may or may not be
//   important.

func (c client) GenerateReport(username, token string, repository []string) string {
	request := &pb.GenerateReportRequest{
		Username:       username,
		Token:          token,
		RepositoryName: repository,
	}
	reply, err := c.GithubStatsClient.GenerateReport(c.Context, request)
	var report string
	for _, r := range reply.Repository {
		report = report + r.GetLatestCommitAuthor() + " "
	}
	if err != nil {
		_ = c.Logger.Log("err", err) // Without an error return parameter, we can't do anything else...
		return ""
	}
	return report
}
