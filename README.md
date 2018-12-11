# goGithubStats
A terminal application that accepts a list of public Github repositories through STDIN pipe and prints out the name, clone URL, date of latest commit and name of latest author for each one.

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
    $ go build
    $ cat ./repolist.txt | ./goGithubStats
    ```

### Docker

    $ docker build -t sh4nnongoh/goGithubStats:latest .
    $ docker run -t sh4nnongoh/goGithubStats:latest
