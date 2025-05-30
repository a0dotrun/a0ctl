package run

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/a0dotrun/a0ctl/internal/appconfig"
	"github.com/a0dotrun/a0ctl/internal/command/builder"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Runs a0 app image."
		long  = "Runs a0 app image based on proved Dockerfile and build config."
	)

	cmd := &cobra.Command{
		Use:   "run [path]",
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

			fmt.Printf("app working directory: %s\n", appPath)

			// Resolve app config from abs path
			// Must run `$ a0 init`
			appConfig, err := appconfig.ResolveAppConfig(appPath)
			if err != nil {
				return fmt.Errorf("failed to resolve app config: %w", err)
			}

			dockerfilePath := builder.ResolveDockerfile(appPath)
			if dockerfilePath == "" {
			}
			if dockerfilePath == "" {
				return fmt.Errorf("no Dockerfile found at %s", appPath)
			}
			// TODO: @sanchitrk
			// create helper method to fetch docker ignores
			dockerIgnorefilePath := ""

			platform := "linux/amd64"
			tag := builder.NewBuildTag(appConfig.AppName, "")

			buildDir, err := builder.BuildOutputDir(appPath)
			if err != nil {
				return err
			}
			fmt.Printf("build output directory: %s\n", buildDir)

			imageOptions := builder.ImageOptions{
				AppName:        appConfig.AppName,
				WorkingDir:     appPath,
				DockerfilePath: dockerfilePath,
				IgnorefilePath: dockerIgnorefilePath,
				Platform:       platform,
				Tag:            tag,
				BuildOutDir:    buildDir,
			}
			ctx := cmd.Context()
			err = ExeImage(ctx, &appConfig, &imageOptions)
			return err
		},
	}
	return cmd
}

func ExeImage(
	ctx context.Context, appConfig *appconfig.Config, imageOptions *builder.ImageOptions) error {
	if err := builder.CheckDockerDaemon(); err != nil {
		return err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("failed to close dagger client: %v", err)
		}
	}(client)

	return nil
}
