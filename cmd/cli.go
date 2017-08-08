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
	envHost  = "SLOPPY_APIHOST"
	envToken = "SLOPPY_APITOKEN"
)

// client is used in each command to handle api requests.
var client *api.Client

func main() {
	stackTrace := false // stackTrace holds the state whether a stack trace is displayed
	defer func() {
		if err := recover(); err != nil {
			printError("Error executing CLI: %s", err)
			if stackTrace {
				debug.PrintStack()
			}
			os.Exit(1)
		}
	}()

	// Shortcut --version, -v to show version command.
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
		if arg == "--help" {
			args = append([]string{"--help"}, args...)
		}
		if arg == "--debug" {
			stackTrace = true
			args = append(args[:i], args[i+1:]...)
		}
	}

	// Update mechanism
	update := make(chan struct{}, 1)
	if len(args) > 0 && args[0] == "version" {
		update <- struct{}{}
	} else {
		go func() {
			if ok, output := checkVersion(); ok {
				fmt.Fprint(os.Stderr, output)
			}
			update <- struct{}{}
		}()
	}

	client = api.NewClient()
	client.SetUserAgent(userAgent())

	if token, ok := os.LookupEnv(envToken); ok {
		client.SetAccessToken(token)
	} else {
		fatal("Missing access token")
	}

	if host, ok := os.LookupEnv(envHost); ok {
		err := client.SetBaseURL(host)
		if err != nil {
			fatal("Error setting base url to %q", host)
		}
		println("Setting client base url to:", host)
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: BasicHelpFunc("sloppy"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fatal("error executing CLI: %s\n", err.Error())
	}

	<-update // wait for update goroutine
	os.Exit(exitCode)
}

func fatal(msg string, args ...interface{}) {
	printError(msg, args)
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
