package auth

import (
	"fmt"

	"github.com/a0dotrun/a0ctl/internal/api"
	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/spf13/cobra"
)

func newWhoAMI() *cobra.Command {
	const (
		use   = "whoami"
		short = "Show the current logged in user or token user."
	)
	cmd := &cobra.Command{
		Use:               use,
		Short:             short,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cli.NoFilesArg,
		RunE:              whoAmI,
	}
	return cmd
}

func whoAmI(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true
	client, err := api.AuthedClient()
	if err != nil {
		return err
	}

	user, err := client.Users.GetUser()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", user.Username)
	return nil
}
