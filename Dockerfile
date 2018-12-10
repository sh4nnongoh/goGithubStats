FROM golang:1.9 as build-base
RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/github.com/sh4nnongoh/goGithubStats
WORKDIR /go/src/github.com/sh4nnongoh/goGithubStats
RUN dep ensure -v && go build

FROM golang:latest
COPY --from=build-base /go/src/github.com/sh4nnongoh/goGithubStats/goGithubStats /usr/local/bin/goGithubStats
CMD echo 'sh4nnongoh/goGithubStats' | goGithubStats