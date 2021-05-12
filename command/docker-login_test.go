package command_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/sloppyio/cli/command"
	"github.com/sloppyio/cli/ui"
)

func TestDockerLoginCommand_implements(t *testing.T) {
	c := &command.DockerLoginCommand{}

	if !strings.Contains(c.Help(), "Usage: sloppy docker-login") {
		t.Errorf("Help = %s", c.Help())
	}

	if !strings.Contains(c.Synopsis(), "") {
		t.Errorf("Synopsis = %s", c.Synopsis())
	}
}

func TestDockerLoginCommand(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	registryCredentials := &mockRegistryCredentialsEndpoint{}
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{InputReader: inR}}
	c := &command.DockerLoginCommand{UI: mockUI, RegistryCredentials: registryCredentials}

	// Create dummy file
	file := createTempFile(t, "docker", "aWxvcmVtOmlwc3Vt")
	defer os.Remove(file.Name())

	args := []string{file.Name()}

	go fmt.Fprintf(inW, "y\n")

	testCodeAndOutput(t, mockUI, c.Run(args), 0, "to our service. You can now launch apps from your private repositories.")
}

func TestDockerLoginCommand_failed(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	registryCredentials := &mockRegistryCredentialsEndpoint{}
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{InputReader: inR}}
	c := &command.DockerLoginCommand{UI: mockUI, RegistryCredentials: registryCredentials}

	// Create dummy file
	file := createTempFile(t, "docker", "")
	defer os.Remove(file.Name())

	args := []string{file.Name()}

	go fmt.Fprintf(inW, "y\n")

	testCodeAndOutput(t, mockUI, c.Run(args), 1, "Unable to upload docker credentials")
}

func TestDockerLoginCommand_abort(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	mockUI := &ui.MockUI{MockUi: &cli.MockUi{InputReader: inR}}
	c := &command.DockerLoginCommand{UI: mockUI}

	// Create dummy file
	file := createTempFile(t, "docker", "")
	defer os.Remove(file.Name())

	args := []string{file.Name()}

	go fmt.Fprintf(inW, "n\n")

	testCodeAndOutput(t, mockUI, c.Run(args), 0, "")
}

func TestDockerLoginCommand_noDockerConfig(t *testing.T) {
	mockUI := &ui.MockUI{MockUi: &cli.MockUi{}}
	c := &command.DockerLoginCommand{UI: mockUI}

	args := []string{"noDockerConfig.json"}
	testCodeAndOutput(t, mockUI, c.Run(args), 1, "doesn't exist.")
}

// createTempFile creates a dummy file for testing purpose.
func createTempFile(t *testing.T, name, auth string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), name)
	if err != nil {
		t.Fatal("Couldn't create temp file!")
	}
	defer file.Close()

	if len(auth) > 0 {
		file.Write([]byte(`{"auths":{"https://index.docker.io/v1/":{"auth":"` + auth + `"}}}`))
	} else {
		file.Write([]byte(`{"auths":{}}}`))
	}
	file.Sync()

	return file
}
