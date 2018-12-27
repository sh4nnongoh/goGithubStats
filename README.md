# goGithubStats
A Go app that accepts a list of public Github repositories and prints out the name, clone URL, date of latest commit and name of latest author for each one.

## Input
- Read plain text list of repositories from stdin
- One repo per input line, format: $orgname/$repo , e.g. kubernetes/charts
## Output
- One line per input repo in CSV format plus one header line to stdout

## Building

### Local

1. Clone goGithubStats

    ```
    $ go get github.com/sh4nnongoh/goGithubStats
    ```

1. Build goGithubStats

    ```
    $ cd $GOPATH/src/github.com/sh4nnongoh/goGithubStats
    $ go build ./src/examples/*/main.go
    ```

### Docker

    ```
    $ cd $GOPATH/src/github.com/sh4nnongoh/goGithubStats
    $ docker build -t sh4nnongoh/goGithubStats .
    $ docker run sh4nnongoh/goGithubStats
    ```