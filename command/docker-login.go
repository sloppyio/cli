package command

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sloppyio/cli/pkg/api"
	"github.com/sloppyio/cli/pkg/dockerconfig"
	"github.com/sloppyio/cli/ui"
)

// DockerLoginCommand is a Command implementation that uploads the docker
// credentials.
type DockerLoginCommand struct {
	UI                  ui.UI
	RegistryCredentials api.RegistryCredentialsUploader
}

// Help should return long-form help text.
func (c *DockerLoginCommand) Help() string {
	helpText := `
Usage: sloppy docker-login [FILENAME]

  Uploads docker credentials in order to access private repositories.

Examples:
  sloppy docker-login
  sloppy docker-login ~/shared/docker/config
`
	return strings.TrimSpace(helpText)
}

// Run should run the actual command with the given CLI instance and
// command-line args.
func (c *DockerLoginCommand) Run(args []string) int {
	dockerConfig := filepath.Join(c.getHomeDir(), ".docker", "config.json")

	if len(args) > 0 {
		dockerConfig = args[0]
	}

	relDockerConfig := dockerConfig
	if rel, err := filepath.Rel(c.getHomeDir(), dockerConfig); err == nil {
		relDockerConfig = filepath.ToSlash("~/" + rel)
	}

	file, err := os.Open(dockerConfig)
	if err != nil {
		if os.IsNotExist(err) {
			c.UI.Error(fmt.Sprintf("%s doesn't exist.", relDockerConfig))
			c.UI.Output("Run 'docker login' to create this file.")
		} else {
			c.UI.Error(err.Error())
		}
		return 1
	}
	defer file.Close()

	c.UI.Warn(fmt.Sprintf("This command will send the content of %s to our service, to give you access to private repos.", relDockerConfig))
	confirm, err := c.UI.Ask("Are you sure you want to continue? (y/n)")
	if err != nil {
		return 1
	}
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		c.UI.Output("Abort.")
		return 0
	}

	transformed, err := dockerconfig.Transform(file)
	if err != nil {
		c.UI.ErrorAPI(err)
		return 1
	}

	if _, _, err := c.RegistryCredentials.Upload(transformed); err != nil {
		c.UI.ErrorAPI(err)
		return 1
	}

	c.UI.Info(fmt.Sprintf("Uploaded %s to our service. You can now launch apps from your private repositories.", relDockerConfig))

	return 0
}

// getHomeDir returns the home directory of the current user with the help of
// environment variables depending on the target operating system.
func (c *DockerLoginCommand) getHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}

	homeDir := os.Getenv(env)
	if homeDir == "" && runtime.GOOS != "windows" {
		if u, err := user.Current(); err == nil {
			return u.HomeDir
		}
	}

	return homeDir
}

// Synopsis should return a one-line, short synopsis of the command.
func (c *DockerLoginCommand) Synopsis() string {
	return "Uploads docker credentials to sloppy.io"
}
