package main

import (
	"fmt"
	"os"
	"runtime"
)

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
	if len(os.Args) < 2 {
		fmt.Println("error: no command provided")
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	handleCommandLine()
}

const (
	appName        = "minyls"
	appDescription = "Go tool to combine MinIO and YOURLS."
	appGithubLink  = "github.com/devusSs/minyls"
)

func printHelp() {
	fmt.Printf("%s - %s\n", appName, appDescription)
	fmt.Println()
	fmt.Println(appGithubLink)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("	minyls <command> <parameters>")
	fmt.Println()
	fmt.Println("Available commands and parameters:")
	fmt.Println("	help")
	fmt.Println("	version")
	fmt.Println("	upload		[filepath] [policy]")
	fmt.Println("	list")
	fmt.Println("	download	[id] [filepath]")
	fmt.Println("	delete		[id]")
	fmt.Println("	clear		[option]")
}

func handleCommandLine() {
	command := os.Args[1]

	switch command {
	case "help":
		printHelp()
		os.Exit(0)
	case "version":
		printVersion()
		os.Exit(0)
	case "upload":
		fmt.Println("upload command, not implemented")
	case "list":
		fmt.Println("list command, not implemented")
	case "download":
		fmt.Println("download command, not implemented")
	case "delete":
		fmt.Println("delete command, not implemented")
	case "clear":
		fmt.Println("clear command, not implemented")
	default:
		fmt.Println("error: unrecognized command:", command)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("%s - %s\n", appName, appDescription)
	fmt.Println()
	fmt.Println(appGithubLink)
	fmt.Println()
	fmt.Printf("Build version:\t\t%s\n", buildVersion)
	fmt.Printf("Build date:\t\t%s\n", buildDate)
	fmt.Printf("Build Git commit:\t%s\n", buildGitCommit)
	fmt.Println()
	fmt.Printf("Build Go version:\t%s\n", runtime.Version())
	fmt.Printf("Build Go OS:\t\t%s\n", runtime.GOOS)
	fmt.Printf("Build Go arch:\t\t%s\n", runtime.GOARCH)
}
