package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sh4nnongoh/goGithubStats/src/githubstatswithextpkg"
)

func main() {
	svc := githubstatswithextpkg.NewService()

	generateReportHandler := httptransport.NewServer(
		githubstatswithextpkg.MakeGenerateReportEndpoint(svc),
		githubstatswithextpkg.DecodeGenerateReportRequest,
		githubstatswithextpkg.EncodeResponse,
	)

	http.Handle("/generateReport", generateReportHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
