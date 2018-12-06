# goGithubStats
An application/script that accepts a list of public Github repositories and prints out the name, clone URL, date of latest commit and name of latest author for each one.

## Input
● Read plain text list of repositories from stdin
● One repo per input line, format: $orgname/$repo , e.g. kubernetes/charts
● Other parameters/env vars as needed, should be documented
## Output
● One line per input repo in CSV or TSV format plus one header line to stdout