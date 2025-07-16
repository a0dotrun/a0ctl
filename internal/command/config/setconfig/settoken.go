package setconfig

import (
	"errors"
	"fmt"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/settings"

	"github.com/a0dotrun/a0ctl/internal/api"
	"github.com/spf13/cobra"
)

func newSetToken() *cobra.Command {
	const (
		use   = "token <jwt>"
		short = "Configure the token used by a0ctl"
	)
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(
			cmd *cobra.Command, args []string, toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: setToken,
	}
	return cmd
}

func setToken(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	config, err := settings.ReadSettings()
	if err != nil {
		return fmt.Errorf("failed to read settings: %w", err)
	}

	token := args[0]
	if !api.IsJWTTokenValid(token) {
		return errors.New("invalid token")
	}

	config.SetToken(token)
	if err := settings.TryToPersistChanges(); err != nil {
		return fmt.Errorf("%w\nIf the issue persists, set your token to the %s environment variable instead", err, cli.Emph(settings.EnvAccessToken))
	}
	fmt.Println("Token set succesfully.")
	return nil
}
