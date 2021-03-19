package command

import (
	"flag"
	"strings"

	"github.com/sloppyio/cli/pkg/api"
	"github.com/sloppyio/cli/ui"
)

// LogsCommand is a Command implementation that is used to fetch standard
// output and standard error streams for a specific application, an entire
// project or service.
type LogsCommand struct {
	UI       ui.UI
	Projects api.ProjectsLogger
	Services api.ServicesLogger
	Apps     api.AppsLogger
}

// Help should return long-form help text.
func (c *LogsCommand) Help() string {
	helpText := `
Usage: sloppy logs [OPTIONS] PROJECT[/SERVICE[/APP]]

  Fetches the logs of the given project, service or app

Options:

	-n                  Number of lines to show from the end of the logs
	--from              Filters the result by start date. Date format is YYYY-MM-DD
	--to                Filters the result by end date. Date format is YYYY-MM-DD


Examples:

	sloppy logs -n 10 letschat
	sloppy logs --from 2018-01-13 --to 2018-01-15 letschat
  sloppy logs letschat/frontend/apache
`
	return strings.TrimSpace(helpText)
}

// Run should run the actual command with the given CLI instance and
// command-line args.
func (c *LogsCommand) Run(args []string) int {
	var lines int
	var fromDate string
	var toDate string

	cmdFlags := newFlagSet("logs", flag.ContinueOnError)
	cmdFlags.IntVar(&lines, "n", 0, "")
	cmdFlags.StringVar(&fromDate, "from", "", "")
	cmdFlags.StringVar(&toDate, "to", "", "")

	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		c.UI.Output("See 'sloppy logs --help'.")
		return 1
	}

	if code := c.UI.ErrorNoFlagAfterArg(cmdFlags.Args()); code == 1 {
		return code
	}

	if cmdFlags.NArg() < 1 {
		return c.UI.ErrorNotEnoughArgs("logs", "", 1)
	}

	parts := strings.Split(strings.Trim(cmdFlags.Arg(0), "/"), "/")

	var logs <-chan api.LogEntry
	var errors <-chan error

	switch len(parts) {
	case 1:
		logs, errors = c.Projects.GetLogs(parts[0], lines, fromDate, toDate)
	case 2:
		logs, errors = c.Services.GetLogs(parts[0], parts[1], lines, fromDate, toDate)
	case 3:
		logs, errors = c.Apps.GetLogs(parts[0], parts[1], parts[2], lines, fromDate, toDate)
	default:
		return c.UI.ErrorInvalidAppPath(cmdFlags.Arg(0))
	}

	for {
		select {
		case err := <-errors:
			if err != nil {
				c.UI.ErrorAPI(err)
			}
			return 1
		case entry, ok := <-logs:
			if !ok {
				return 0
			}
			c.UI.Output(entry.String())
		}
	}
}

// Synopsis should return a one-line, short synopsis of the command.
func (c *LogsCommand) Synopsis() string {
	return "Fetch the logs of a project, service or app"
}
