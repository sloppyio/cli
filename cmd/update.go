package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/sloppyio/cli/command"
)

var updateURL = "https://files.sloppy.io/version.txt"
var updateFilename = "sloppy_updateNotifier"

const updateFormat = time.RFC3339

func checkVersion() (bool, string) {
	if ok, _ := updateFile(true); ok {
		return false, ""
	}

	deployedVersion, err := getDeployedVersion()
	if err != nil {
		return false, "Could not look for newer versions"
	}

	if compareVersion(deployedVersion, command.Version) == 1 {
		// Truncate file to force update request
		updateFile(false)
		return true, fmt.Sprintf(`A newer version of this runtime is available
Server has version: %s
User has version: %s

Check https://sloppy.io/knowledge-base/install-cli/ for install/update instructions
and https://github.com/sloppyio/cli/blob/master/CHANGELOG.md for the changelog`, deployedVersion, command.Version)
	}

	return false, ""
}

func updateFile(update bool) (bool, error) {
	file, err := os.OpenFile(filepath.Join(os.TempDir(), updateFilename), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return false, err
	}

	defer func() {
		file.Truncate(0)
		if update {
			fmt.Fprint(file, time.Now().Format(updateFormat))
		}
		file.Sync()
		file.Close()
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return false, err
	}

	if date := string(data); date != "" {
		lastCheckDate, err := time.Parse(updateFormat, date)
		if err != nil {
			return false, nil
		}
		if lastCheckDate.After(time.Now().Truncate(24 * time.Hour)) {
			return true, nil
		}
	}

	return false, nil
}

func getDeployedVersion() (string, error) {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Get(updateURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func compareVersion(a, b string) int {
	aVersion, err := version.NewVersion(a)
	if err != nil {
		return -1
	}
	bVersion, err := version.NewVersion(b)
	if err != nil {
		return 1
	}
	return aVersion.Compare(bVersion)
}
