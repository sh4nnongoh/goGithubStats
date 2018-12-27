package main

import (
	"log"
	"net/http"

	"github.com/sh4nnongoh/goGithubStats/src/githubstats"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	svc := githubstats.NewGithubReportService()

	generateReportHandler := httptransport.NewServer(
		githubstats.MakeGenerateReportEndpoint(svc),
		githubstats.DecodeGenerateReportRequest,
		githubstats.EncodeResponse,
	)

	http.Handle("/generateReport", generateReportHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
