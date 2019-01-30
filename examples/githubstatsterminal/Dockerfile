FROM golang:1.9 as build-base
RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/github.com/sh4nnongoh/githubstatsterminal
WORKDIR /go/src/github.com/sh4nnongoh/githubstatsterminal
RUN dep ensure -v && go build 

FROM golang:latest
COPY --from=build-base /go/src/github.com/sh4nnongoh/githubstatsterminal/githubstatsterminal /usr/local/bin/githubstatsterminal
CMD echo 'sh4nnongoh/goGithubStats' | githubstatsterminal