package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/sh4nnongoh/goGithubStats/src/githubstatscmd"
)

func main() {
	var (
		username = flag.String("username", "", "Github username")
		token    = flag.String("token", "", "Github personel access token")
	)
	flag.Parse()

	// Obtain input stream of repositories through stdin
	var repositoryList []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		repositoryList = append(repositoryList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	githubReport := githubstatscmd.NewGithubReport(*username, *token, repositoryList)
	githubReport.PrintRepositoryDetails()
}
