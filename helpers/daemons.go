package helpers

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func CheckDockerDaemon() error {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "daemon") || strings.Contains(stderr, "connect") {
				return fmt.Errorf(
					"docker daemon is not running. Please start Docker and try again")
			}
		}
		return fmt.Errorf(
			"docker is not available or not installed. Please install Docker and ensure it's running")
	}

	if len(strings.TrimSpace(string(output))) == 0 {
		return fmt.Errorf("docker daemon is not running. Please start Docker and try again")
	}

	return nil
}
