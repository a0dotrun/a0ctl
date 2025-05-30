package builder

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/a0dotrun/a0ctl/internal/appconfig"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type ImageOptions struct {
	AppName        string
	WorkingDir     string
	DockerfilePath string
	IgnorefilePath string
	Tag            string
	Label          map[string]string
	Platform       string
}

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Build a0 app image."
		long  = "Build a0 app image based on proved Dockerfile and build config."
	)

	return &cobra.Command{
		Use:   "build [path]",
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			buildPath := "."
			if len(args) > 0 {
				buildPath = args[0]
			}

			appPath, err := filepath.Abs(buildPath)
			if err != nil {
				return fmt.Errorf("failed to resolve app path: %w", err)
			}

			fmt.Printf("Build: %s\n", appPath)

			dockerfilePath := ResolveDockerfile(appPath)
			if dockerfilePath == "" {
				return fmt.Errorf("no Dockerfile found at %s", appPath)
			}
			// TODO: @sanchitrk
			// create helper method to fetch docker ignores
			dockerIgnorefilePath := ""

			// TODO: @sanchitrk
			// read from config file or global appconfig.Config store
			name := "foobar"
			region := "ap-south-1"
			platform := "linux/amd64"
			tag := NewBuildTag(name, "")

			appConfig := appconfig.NewConfig(name, region)

			imageOptions := ImageOptions{
				AppName:        name,
				WorkingDir:     appPath,
				DockerfilePath: dockerfilePath,
				IgnorefilePath: dockerIgnorefilePath,
				Platform:       platform,
				Tag:            tag,
			}

			ctx := cmd.Context()
			img, err := DetermineImage(ctx, &appConfig, &imageOptions)
			if err != nil {
				return err
			}

			fmt.Printf("built image with dagger...")
			fmt.Println(img)

			return nil
		},
	}
}

type DeploymentImage struct {
	ImageUrl string
}

func DetermineImage(
	ctx context.Context, appConfig *appconfig.Config, imageOptions *ImageOptions,
) (img *DeploymentImage, err error) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return &DeploymentImage{}, err
	}

	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("failed to close dagger client: %v", err)
		}
	}(client)

	imageUrl, err := client.Host().Directory(imageOptions.WorkingDir).
		DockerBuild().
		Publish(ctx, "ttl.sh/"+imageOptions.Tag)

	if err != nil {
		return &DeploymentImage{}, err
	}
	deploymentImage := &DeploymentImage{
		ImageUrl: imageUrl,
	}
	return deploymentImage, nil
}

func NewBuildTag(appName string, label string) string {
	if label == "" {
		label = fmt.Sprintf("build-%s", ulid.Make())
	}
	return fmt.Sprintf("%s:%s", appName, label)
}
