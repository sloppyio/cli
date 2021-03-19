package command

import (
	"bytes"
	"fmt"

	"github.com/sloppyio/cli/ui"
)

var (
	// GitCommit is the commit hash of the current build.
	GitCommit string

	// Version number that is being run at the moment.
	Version string

	// VersionPreRelease marks the version as pre-release. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPreRelease string
)

// VersionCommand is a Command implementation that prints the version.
type VersionCommand struct {
	UI ui.UI
}

// Help should return long-form help text.
func (c *VersionCommand) Help() string {
	return ""
}

// Run should run the actual command with the given CLI instance and
// command-line args.
func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "Sloppy %s", Version)
	if VersionPreRelease != "" {
		fmt.Fprintf(&versionString, ".%s", VersionPreRelease)

		if GitCommit != "" {
			fmt.Fprintf(&versionString, " (%s)", GitCommit)
		}
	}

	c.UI.Output(versionString.String())

	return 0
}

// Synopsis should return a one-line, short synopsis of the command.
func (c *VersionCommand) Synopsis() string {
	return "Prints the sloppy version"
}
