package builder

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
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
			fmt.Printf("I am here to buil your app image")
			return nil
		},
	}
}

type AppConfig struct{}

type DeploymentImage struct{}

func DetermineImage(ctx context.Context, appConfig *AppConfig) (img *DeploymentImage, err error) {
	return &DeploymentImage{}, nil
}
