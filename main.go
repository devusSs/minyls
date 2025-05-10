package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/devusSs/minyls/internal/cli"
	"github.com/devusSs/minyls/internal/log"
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

const minArgs = 2

func main() {
	if len(os.Args) < minArgs {
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

// logging may be used here for cli commands
// since each function in the cli package
// calls cli.initialize first which sets up logging.
func handleCommandLine() {
	command := os.Args[1]

	switch command {
	case "help":
		printHelp()
	case "version":
		printVersion()
	case "upload":
		err := cli.Upload()
		if err != nil {
			log.Log().Err(err).Str("func", "handleCommandLine").Msg("upload failed")
			os.Exit(1)
		}
	case "list":
		err := cli.List()
		if err != nil {
			log.Log().Err(err).Str("func", "handleCommandLine").Msg("list failed")
			os.Exit(1)
		}
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
