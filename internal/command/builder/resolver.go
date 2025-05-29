package builder

import (
	"github.com/a0dotrun/a0ctl/helpers"
	"path/filepath"
)

// ResolveDockerfile - Resolves the location of the dockerfile, allowing for upper and lowercase naming
func ResolveDockerfile(cwd string) string {
	dockerfilePath := filepath.Join(cwd, "Dockerfile")
	if helpers.FileExists(dockerfilePath) {
		return dockerfilePath
	}
	dockerfilePath = filepath.Join(cwd, "dockerfile")
	if helpers.FileExists(dockerfilePath) {
		return dockerfilePath
	}
	return ""
}
