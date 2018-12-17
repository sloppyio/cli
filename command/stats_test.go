package command_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mitchellh/cli"

	"github.com/sloppyio/cli/command"
	"github.com/sloppyio/cli/internal/test"
	"github.com/sloppyio/cli/ui"
)

func TestStatsCommand_implements(t *testing.T) {
	c := &command.StatsCommand{}

	if !strings.Contains(c.Help(), "Usage: sloppy stats") {
		t.Errorf("Help = %s", c.Help())
	}

	if !strings.Contains(c.Synopsis(), "Display metrics") {
		t.Errorf("Synopsis = %s", c.Synopsis())
	}
}

func TestStatsCommand(t *testing.T) {
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{}}
	projects := &mockProjectsEndpoint{}
	apps := &mockAppsEndpoint{}
	c := &command.StatsCommand{UI: mockUI, Projects: projects, Apps: apps}

	args := []string{
		"letschat",
	}
	testCodeAndOutput(t, mockUI, c.Run(args), 1, "")
	out := mockUI.OutputWriter.String()
	if !strings.Contains(out, "No apps running") {
		t.Errorf("Output = %s", out)
	}
}

// This test introduces a native client server connection and
// loading content from testdata folder instead of mocking all of this.
func TestStatsCommand_withAllFlag(t *testing.T) {
	helper := test.NewHelper(t)
	projectResponse := helper.LoadFile("letschat_response_project.json")
	statsResponse := helper.LoadFile("letschat_response_stats.json")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/apps/letschat") {
			n, err := w.Write(projectResponse)
			if err != nil || n != len(projectResponse) {
				t.Error(err)
			}
		} else if strings.HasSuffix(r.URL.Path, "/apps/letschat/services/frontend/apps/node/stats") {
			n, err := w.Write(statsResponse)
			if err != nil || n != len(statsResponse) {
				t.Error(err)
			}
		} else {
			http.NotFound(w, r)
		}
	}
	server := helper.NewAPIServer(handler)
	defer server.Close()

	client := helper.NewClient(server.Listener.Addr())
	client.SetAccessToken("gimmeAccess")

	mockUI := &ui.MockUI{MockUi: &cli.MockUi{}}
	c := &command.StatsCommand{UI: mockUI, Projects: client.Projects, Apps: client.Apps}
	args := []string{
		"--all",
		"letschat",
	}

	// FIXME - table head tab differs between tests
	//	want := `CONTAINER 		 CPU % 	 MEM / LIMIT 		 MEM % 	 NET I/O Extern 	 NET I/O Intern
	//frontend/node-59f7ed 	 0.0% 	 128 MiB / 1024 MiB 	 12.5% 	 7.71 B / 0 B 	 18.1 B / 0 B`

	if exitCode := c.Run(args); exitCode != 0 {
		t.Error(mockUI.ErrorWriter.String())
	}
	//else if diff := cmp.Diff(mockUI.OutputWriter.String(), want); diff != "" {
	//	t.Errorf("Result differs: (-got +want)\n%s", diff)
	//}
}

func TestStatsCommand_notEnoughArgs(t *testing.T) {
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{}}
	c := &command.StatsCommand{UI: mockUI}

	args := []string{}
	testCodeAndOutput(t, mockUI, c.Run(args), 1, "minimum of 1 argument")
}

func TestStatsCommand_notFound(t *testing.T) {
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{}}
	projects := &mockProjectsEndpoint{}
	c := &command.StatsCommand{UI: mockUI, Projects: projects}

	args := []string{
		"abc",
	}

	testCodeAndOutput(t, mockUI, c.Run(args), 1, "not be found")
}

func TestStatsCommand_invalidProjectPath(t *testing.T) {
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{}}
	projects := &mockProjectsEndpoint{}
	c := &command.StatsCommand{UI: mockUI, Projects: projects}

	args := []string{
		"abc/def",
	}

	testCodeAndOutput(t, mockUI, c.Run(args), 1, "invalid project path")
}
