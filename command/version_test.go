package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/sloppyio/cli/ui"
)

func TestVersionCommand_implements(t *testing.T) {
	c := &VersionCommand{}

	if !strings.Contains(c.Help(), "") {
		t.Errorf("Help = %s", c.Help())
	}

	if !strings.Contains(c.Synopsis(), "sloppy version") {
		t.Errorf("Synopsis = %s", c.Synopsis())
	}
}

func TestVersionCommand(t *testing.T) {
	mockUI := &ui.MockUI{MockUi: new(cli.MockUi)}
	Version, VersionPreRelease, GitCommit = "0.0.1", "dev", "1b33f1"
	c := &VersionCommand{
		CheckVersion: func() (bool, string) {
			return false, ""
		},
		UI: mockUI,
	}

	testCodeAndOutput(t, mockUI, c.Run(nil), 0, "Sloppy 0.0.1.dev (1b33f1)")
}
