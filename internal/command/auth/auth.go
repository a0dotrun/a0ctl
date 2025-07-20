// Package auth provides the authentication related commands.
package auth

import (
	"fmt"
	"os"

	"github.com/a0dotrun/a0ctl/internal/api"
	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/flags"
	"github.com/a0dotrun/a0ctl/internal/settings"
	"github.com/spf13/cobra"
)

// authURLPath is the path to the authentication endpoint.
const authURLPath = "/auth/cli"

func New() *cobra.Command {
	const (
		short = "Manage authentication"
		long  = "Authenticate with a0 (and logout if you need to)."
	)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: short,
		Long:  long,
	}

	loginCmd := newLogin()
	flags.AddHeadless(loginCmd)
	cmd.AddCommand(newWhoAMI(), loginCmd)

	return cmd
}

func checkEnvAuth(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	token := os.Getenv(settings.EnvAccessToken)
	if token != "" {
		return fmt.Errorf("a token is set in the %q environment variable, please unset it before running %s", settings.EnvAccessToken, cli.Emph(cmd.CommandPath()))
	}
	return nil
}

func validateToken(token string) (string, error) {
	client, err := api.MakeClient(token)
	if err != nil {
		return "", fmt.Errorf("could not create client to validate token: %w", err)
	}

	user, err := client.Users.GetUser()
	if err != nil {
		return "", fmt.Errorf("could not validate token: %w", err)
	}

	return user.Username, nil
}
