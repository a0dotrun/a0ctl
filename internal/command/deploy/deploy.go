// Package deploy provides the deploy command and helpers for the CLI.
package deploy

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/a0dotrun/a0ctl/internal/api"
	"github.com/a0dotrun/a0ctl/internal/command/deploy/utils"
	"github.com/a0dotrun/a0ctl/internal/settings"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	const (
		short = "Deploy code to the server"
		long  = "Build the docker image (using Dockerfile) and deploys it to the server"
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: short,
		Long:  long,
		RunE:  deploy,
	}

	return cmd
}

func deploy(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	fmt.Printf("Deploying to the server...\n")
	settings, err := settings.ReadSettings()
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	artifactDir := settings.GetLocalConfigDir()
	artifactPath := filepath.Join(artifactDir, "artifact.tar.gz")

	client, err := api.AuthedClient()
	if err != nil {
		return fmt.Errorf("Error creating client")
	}

	err = utils.CreateArtifact("./examples/app", artifactPath)
	if err != nil {
		panic(err)
	}

	artifactResp, err := utils.CallArtifact(client)
	if err != nil {
		return fmt.Errorf("failed to get upload URL: %w", err)
	}
	log.Printf("Upload URL: %s", artifactResp.UploadURL)
	log.Printf("Expires: %d", artifactResp.Expires)
	log.Printf("Artifact Key: %s", artifactResp.ArtifactKey)

	err = utils.UploadArtifact(artifactPath, artifactResp.UploadURL)
	if err != nil {
		return fmt.Errorf("failed to upload artifact: %w", err)
	}

	deployResp, err := utils.CallDeploy(client, artifactResp.ArtifactKey)
	if err != nil {
		return fmt.Errorf("Failed to deploy the artifact..., %v", err)
	}

	log.Printf("DeploymentId: %s", deployResp.DeploymentId)
	log.Printf("Status: %v", deployResp.Status)
	log.Printf("CreatedAt: %v", deployResp.CreatedAt)

	return nil
}
