package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/sloppyio/cli/command"
	"github.com/sloppyio/cli/pkg/api"
)

const (
	envAPIURL = "SLOPPY_API_URL"
	envToken  = "SLOPPY_APITOKEN"
)

// client is used in each command to handle api requests.
var client *api.Client

func main() {
	stackTrace := false // stackTrace holds the state whether a stack trace is displayed
	defer func() {
		if err := recover(); err != nil {
			printError("%s\nFor help, please visit https://kb.sloppy.io/features#cli-command-reference", err)
			if stackTrace {
				debug.PrintStack()
			}
			os.Exit(1)
		}
	}()

	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--help" {
			args = append([]string{"--help"}, args...)
		}
		if arg == "--debug" {
			stackTrace = true
			args = append(args[:i], args[i+1:]...)
		}
	}

	client = api.NewClient()
	client.SetUserAgent(userAgent())

	if token, ok := os.LookupEnv(envToken); ok {
		client.SetAccessToken(token)
	} else if len(args) != 0 && args[0] != "--help" {
		fatal("Missing %s, please login by exporting your token https://admin.sloppy.io/account/tokens", envToken)
	}

	if apiURL, ok := os.LookupEnv(envAPIURL); ok {
		err := client.SetBaseURL(apiURL)
		if err != nil {
			fatal("Error setting base url to %q", apiURL)
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: BasicHelpFunc("sloppy"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fatal("%s\nFor help, please visit https://kb.sloppy.io/features#cli-command-reference", err.Error())
	}

	os.Exit(exitCode)
}

func fatal(msg string, args ...interface{}) {
	printError(msg, args...)
	os.Exit(1)
}

func printError(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	w := os.Stderr
	if runtime.GOOS == "windows" {
		fmt.Fprint(w, message)
	} else {
		fmt.Fprintf(w, "\033[0;31m%s\033[0m\n", message)
	}
}

func userAgent() string {
	sys := strings.Title(strings.Replace(runtime.GOOS, "darwin", "macintosh", -1))
	agent := fmt.Sprintf("sloppy-cli/%s", command.Version)
	if command.VersionPreRelease != "" {
		agent += fmt.Sprintf(".%s", command.VersionPreRelease)
	}
	agent += fmt.Sprintf(" (%s) go/%s", sys, runtime.Version()[2:])
	return agent
}
