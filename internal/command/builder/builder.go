package builder

import (
	"context"
	"dagger.io/dagger"
	"errors"
	"fmt"
	"github.com/a0dotrun/a0ctl/internal/appconfig"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
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
	Label          map[string]string
	Platform       string
	BuildOutDir    string
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
			publish, _ := cmd.Flags().GetBool("publish")
			version, _ := cmd.Flags().GetString("version")
			if version == "" {
				version = "latest"
			}

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

			platform := "linux/amd64"
			tag := NewBuildTag(appConfig.AppName, "", version)

			buildDir, err := BuildOutputDir(appPath)
			if err != nil {
				return err
			}
			fmt.Printf("build output directory: %s\n", buildDir)

			imageOptions := ImageOptions{
				AppName:        appConfig.AppName,
				WorkingDir:     appPath,
				DockerfilePath: dockerfilePath,
				IgnorefilePath: dockerIgnorefilePath,
				Platform:       platform,
				Tag:            tag,
				BuildOutDir:    buildDir,
			}
			ctx := cmd.Context()
			img, err := DetermineImage(ctx, &appConfig, &imageOptions, publish)
			if err != nil {
				return err
			}

			fmt.Printf("built and published image: %s\n", img.ImageUrl)
			return nil
		},
	}

	cmd.Flags().Bool("publish", false, "Publish the built image to registry")

	return cmd
}

type DeploymentImage struct {
	ImageUrl string
}

func DetermineImage(
	ctx context.Context, appConfig *appconfig.Config, imageOptions *ImageOptions, publish bool,
) (img *DeploymentImage, err error) {
	if err := CheckDockerDaemon(); err != nil {
		return &DeploymentImage{}, err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return &DeploymentImage{}, fmt.Errorf("failed to connect to Dagger: %w", err)
	}
	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("failed to close dagger client: %v", err)
		}
	}(client)

	deploymentImage := &DeploymentImage{}

	container := client.Host().Directory(imageOptions.WorkingDir).DockerBuild()

	if publish {
		imageUrl, err := container.Publish(ctx, "ttl.sh/"+imageOptions.Tag)
		if err != nil {
			return &DeploymentImage{}, fmt.Errorf("failed to build and publish image: %w", err)
		}
		deploymentImage.ImageUrl = imageUrl
		return deploymentImage, nil
	}

	imageUrl, err := container.Export(ctx, fmt.Sprintf("%s/%s", imageOptions.BuildOutDir, imageOptions.Tag))
	if err != nil {
		return &DeploymentImage{}, fmt.Errorf("failed to export image to local Docker: %w", err)
	}

	deploymentImage.ImageUrl = imageUrl
	return deploymentImage, nil
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

func NewBuildTag(appName string, label, version string) string {
	if label == "" {
		label = fmt.Sprintf("build-%s", ulid.Make())
	}
	return fmt.Sprintf("%s:%s", appName, label)
}

func BuildOutputDir(appPath string) (string, error) {
	a0Dir := filepath.Join(appPath, ".a0")
	if _, err := os.Stat(a0Dir); os.IsNotExist(err) {
		return "", fmt.Errorf(".a0 directory not found at %s", appPath)
	}

	buildsDir := filepath.Join(a0Dir, "builds")
	if err := os.MkdirAll(buildsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create builds directory: %w", err)
	}

	return buildsDir, nil
}
