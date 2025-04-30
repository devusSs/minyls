package main

import "fmt"

var (
	buildVersion   string
	buildDate      string
	buildGitCommit string
)

func init() {
	if buildVersion == "" {
		buildVersion = "development"
	}

	if buildDate == "" {
		buildDate = "unknown"
	}

	if buildGitCommit == "" {
		buildGitCommit = "unknown"
	}
}

func main() {
	fmt.Printf("init: version: %s, date: %s, commit: %s\n", buildVersion, buildDate, buildGitCommit)
}
