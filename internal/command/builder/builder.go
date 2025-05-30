package builder

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Build a0 app image."
		long  = "Build a0 app image based on proved Dockerfile and build config."
	)

	return &cobra.Command{
		Use:   "build",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("building image with dagger...")
			ctx := cmd.Context()
			app := AppConfig{}
			img, err := DetermineImage(ctx, &AppConfig{})
			if err != nil {
				return err
			}

			fmt.Printf("built image with dagger...")
			fmt.Println(app)
			fmt.Println(img)

			return nil
		},
	}
}

type AppConfig struct{}

type DeploymentImage struct{}

func DetermineImage(ctx context.Context, appConfig *AppConfig) (img *DeploymentImage, err error) {
	client, err := dagger.Connect(ctx,
		dagger.WithLogOutput(os.Stderr),
	)
	if err != nil {
		return &DeploymentImage{}, err
	}

	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("failed to close dagger client: %v", err)
		}
	}(client)

	return &DeploymentImage{}, nil
}
