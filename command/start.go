package command

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sloppyio/cli/pkg/api"
	"github.com/sloppyio/cli/ui"
)

// StartCommand is a Command implementation that is used to create a project
// along with all its services and applications.
type StartCommand struct {
	UI       ui.UI
	Projects api.ProjectsCreater
}

// Help should return long-form help text.
func (c *StartCommand) Help() string {
	helpText := `
Usage: sloppy start [OPTIONS] FILENAME

  Create a new project. Can also use 'sloppy change' for this.

Options:
  -v, --var=[]     values to set for placeholders
  -p, --project    project name

Examples:

  sloppy start sloppy.json
  sloppy start --var=domain:mydomain.sloppy.zone --var=memory:128 myproject.json
  sloppy start --project=myproject docker-compose.yml
`
	return strings.TrimSpace(helpText)
}

// Run should run the actual command with the given CLI instance and
// command-line args.
func (c *StartCommand) Run(args []string) int {
	var vars StringMap
	var projectName string
	cmdFlags := newFlagSet("start", flag.ContinueOnError)
	cmdFlags.Var(&vars, "v", "")
	cmdFlags.Var(&vars, "var", "")
	cmdFlags.StringVar(&projectName, "p", "", "")
	cmdFlags.StringVar(&projectName, "project", "", "")

	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		c.UI.Output("See 'sloppy start --help'.")
		return 1
	}

	if code := c.UI.ErrorNoFlagAfterArg(cmdFlags.Args()); code == 1 {
		return code
	}

	if cmdFlags.NArg() < 1 {
		return c.UI.ErrorNotEnoughArgs("start", "", 1)
	}

	file, err := os.Open(cmdFlags.Arg(0))
	if err != nil {
		if os.IsNotExist(err) {
			c.UI.Error(fmt.Sprintf("file '%s' not found.", cmdFlags.Arg(0)))
		} else if os.IsPermission(err) {
			c.UI.Error(fmt.Sprintf("no read permission '%s'.", cmdFlags.Arg(0)))
		} else {
			c.UI.Error(err.Error())
		}
		return 1
	}
	defer file.Close()

	var inputSource io.Reader
	inputSource = file

	ext := filepath.Ext(file.Name())
	if ext == ".yaml" || ext == ".yml" {
		newSource, err := tryDockerCompose(file.Name(), projectName)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Converting docker-compose failed: %s", err))
			return 1
		}

		if newSource != nil {
			inputSource = newSource
		}
	}

	decoder := newDecoder(inputSource, vars)
	var input = new(api.Project)

	switch ext {
	case ".json":
		if err := decoder.DecodeJSON(input); err != nil {
			c.UI.Error(fmt.Sprintf("failed to parse JSON file %s\n%s", file.Name(), err.Error()))
			return 1
		}
	case ".yaml", ".yml":
		if err := decoder.DecodeYAML(input); err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	default:
		c.UI.Error("file extension not supported, must be json or yaml.")
		return 1
	}

	project, _, err := c.Projects.Create(input)
	if err != nil {
		c.UI.ErrorAPI(err)
		return 1
	}

	c.UI.Table("show", project.Services)
	return 0
}

// Synopsis should return a one-line, short synopsis of the command.
func (c *StartCommand) Synopsis() string {
	return "Create a new project"
}
