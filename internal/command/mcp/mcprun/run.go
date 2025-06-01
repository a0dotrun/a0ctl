package mcprun

import (
	"bufio"
	"context"
	"fmt"
	"github.com/a0dotrun/a0ctl/helpers"
	"github.com/a0dotrun/a0ctl/internal/appconfig"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

type RunOptions struct {
	AppName     string
	WorkingDir  string
	Tag         string
	Ports       []string
	Environment []string
	Detach      bool
	Remove      bool
}

// New initializes and returns a new run Command.
func New() *cobra.Command {
	const (
		short = "Run a0 app container."
		long  = "Run a0 app container from the built image using --tag flag."
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

			tag, _ := cmd.Flags().GetString("tag")
			ports, _ := cmd.Flags().GetStringSlice("port")
			env, _ := cmd.Flags().GetStringSlice("env")
			detach, _ := cmd.Flags().GetBool("detach")
			remove, _ := cmd.Flags().GetBool("rm")

			if tag == "" {
				return fmt.Errorf("--tag is required. Please specify the image tag to run (e.g., --tag latest)")
			}

			appPath, err := filepath.Abs(buildPath)
			if err != nil {
				return fmt.Errorf("failed to resolve app path: %w", err)
			}
			// Resolve app config from abs path to get appName
			appConfig, err := appconfig.ResolveAppConfig(appPath)
			if err != nil {
				return fmt.Errorf("failed to resolve app config: %w", err)
			}

			runOptions := RunOptions{
				AppName:     appConfig.AppName,
				WorkingDir:  appPath,
				Tag:         tag,
				Ports:       ports,
				Environment: env,
				Detach:      detach,
				Remove:      remove,
			}

			ctx := cmd.Context()
			containerId, err := RunContainer(ctx, &runOptions)
			if err != nil {
				return err
			}

			if detach {
				fmt.Printf("container started: %s\n", containerId)
			} else {
				fmt.Printf("container finished: %s\n", containerId)
			}

			return nil
		},
	}

	cmd.Flags().String("tag", "", "Tag of the image to run (e.g., latest, v1.0.0)")
	cmd.Flags().StringSliceP("port", "p", []string{}, "Port mappings (e.g., 8080:80, 3000:3000)")
	cmd.Flags().StringSliceP("env", "e", []string{}, "Environment variables (e.g., KEY=value)")
	cmd.Flags().BoolP("detach", "d", false, "Run container in detached mode")
	cmd.Flags().Bool("rm", true, "Automatically remove container when it exits")

	return cmd
}

func RunContainer(ctx context.Context, runOptions *RunOptions) (string, error) {
	if err := helpers.CheckDockerDaemon(); err != nil {
		return "", err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer func(dockerClient *client.Client) {
		err := dockerClient.Close()
		if err != nil {

		}
	}(dockerClient)

	// Check if image exists
	_, _, err = dockerClient.ImageInspectWithRaw(ctx, fmt.Sprintf("%s:%s", runOptions.AppName, runOptions.Tag))
	if err != nil {
		return "", fmt.Errorf(
			"image %s not found. Please build the image first using 'a0 build'", runOptions.Tag)
	}

	// Check if container with same name exists and remove it
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+runOptions.AppName {
				fmt.Printf("Removing existing container: %s\n", runOptions.AppName)
				err := dockerClient.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true})
				if err != nil {
					return "", fmt.Errorf("failed to remove existing container: %w", err)
				}
				break
			}
		}
	}

	// Parse port mappings
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, portMapping := range runOptions.Ports {
		parts := strings.Split(portMapping, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid port mapping format: %s. Use format hostPort:containerPort", portMapping)
		}

		hostPort := parts[0]
		containerPort := parts[1]

		port, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return "", fmt.Errorf("invalid container port: %s", containerPort)
		}

		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostPort: hostPort,
			},
		}
	}

	// qualified name of the image
	image := runOptions.AppName + ":" + runOptions.Tag

	// Create container configuration
	config := &container.Config{
		Image:        image,
		Env:          runOptions.Environment,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		AutoRemove:   runOptions.Remove,
	}

	// Create container with appName as container name
	resp, err := dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, runOptions.AppName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Flush stdout/stderr before streaming container logs
	os.Stdout.Sync()
	os.Stderr.Sync()

	// Ensure container cleanup on exit/interrupt
	defer func() {
		if !runOptions.Detach {
			fmt.Printf("\nStopping and removing container: %s\n", runOptions.AppName)
			// Stop the container first
			timeout := 3
			err := dockerClient.ContainerStop(context.Background(), resp.ID, container.StopOptions{Timeout: &timeout})
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "failed to stop container, force stopping...")
			}
			// Force stop the container
			_ = dockerClient.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
			if err != nil {
				return
			}
		}
	}()

	if runOptions.Detach {
		return resp.ID, nil
	}

	// Stream container logs to stdout/stderr
	logOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}

	logs, err := dockerClient.ContainerLogs(ctx, resp.ID, logOptions)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	// Stream logs to stdout with prefix in a goroutine
	go func() {
		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			fmt.Printf("[%s]: %s\n", runOptions.AppName, scanner.Text())
		}
	}()

	// Wait for container to finish if not detached
	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container: %w", err)
		}
	case <-statusCh:
	}

	return resp.ID, nil
}
