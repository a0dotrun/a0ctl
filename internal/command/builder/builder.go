package builder

import (
	"context"
	"errors"
	"fmt"
	"github.com/a0dotrun/a0ctl/internal/appconfig"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ImageOptions struct {
	AppName        string
	WorkingDir     string
	DockerfilePath string
	IgnorefilePath string
	Tag            string
	Labels         map[string]string
	Platform       string
}

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Build a0 app image."
		long  = "Build a0 app image based on proved Dockerfile and build config."
	)

	cmd := &cobra.Command{
		Use:   "build [path]",
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			buildPath := "."
			if len(args) > 0 {
				buildPath = args[0]
			}
			tag, _ := cmd.Flags().GetString("tag")
			platform, _ := cmd.Flags().GetString("platform")

			appPath, err := filepath.Abs(buildPath)
			if err != nil {
				return fmt.Errorf("failed to resolve app path: %w", err)
			}
			fmt.Printf("app working directory: %s\n", appPath)

			// Resolve app config from abs path
			// Must run `$ a0 init`
			appConfig, err := appconfig.ResolveAppConfig(appPath)
			if err != nil {
				return fmt.Errorf("failed to resolve app config: %w", err)
			}

			dockerfilePath := ResolveDockerfile(appPath)
			if dockerfilePath == "" {
				return fmt.Errorf("no Dockerfile found at %s", appPath)
			}
			// TODO: @sanchitrk
			// create helper method to fetch docker ignores
			dockerIgnorefilePath := ""

			tag = NewBuildTag(appConfig.AppName, tag)

			imageOptions := ImageOptions{
				AppName:        appConfig.AppName,
				WorkingDir:     appPath,
				DockerfilePath: dockerfilePath,
				IgnorefilePath: dockerIgnorefilePath,
				Platform:       platform,
				Tag:            tag,
			}
			ctx := cmd.Context()
			img, err := BuildImage(ctx, &imageOptions)
			if err != nil {
				return err
			}

			fmt.Printf("built image: %s\n", img.ImageUrl)
			return nil
		},
	}

	cmd.Flags().String("tag", "", "Tag for the built image (e.g., latest, v1.0.0)")
	cmd.Flags().String(
		"platform", "", "Platform for the built image (e.g., linux/amd64, linux/arm64)")

	return cmd
}

type DeploymentImage struct {
	ImageUrl string
}

func BuildImage(ctx context.Context, imageOptions *ImageOptions) (img *DeploymentImage, err error) {
	if err := CheckDockerDaemon(); err != nil {
		return &DeploymentImage{}, err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &DeploymentImage{}, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer func(dockerClient *client.Client) {
		err := dockerClient.Close()
		if err != nil {

		}
	}(dockerClient)

	buildContext, err := archive.TarWithOptions(imageOptions.WorkingDir, &archive.TarOptions{})
	if err != nil {
		return &DeploymentImage{}, fmt.Errorf("failed to create build context: %w", err)
	}
	defer func(buildContext io.ReadCloser) {
		err := buildContext.Close()
		if err != nil {

		}
	}(buildContext)

	buildOptions := build.ImageBuildOptions{
		Tags:       []string{imageOptions.Tag},
		Dockerfile: filepath.Base(imageOptions.DockerfilePath),
		Remove:     true,
		Platform:   imageOptions.Platform,
		Labels:     imageOptions.Labels,
	}

	buildResponse, err := dockerClient.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return &DeploymentImage{}, fmt.Errorf("failed to build image: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(buildResponse.Body)

	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		return &DeploymentImage{}, fmt.Errorf("failed to read build output: %w", err)
	}

	return &DeploymentImage{ImageUrl: imageOptions.Tag}, nil
}

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

func NewBuildTag(appName, tag string) string {
	if tag == "" {
		tag = fmt.Sprintf("build-%s", ulid.Make())
	}
	return fmt.Sprintf("%s:%s", appName, tag)
}
