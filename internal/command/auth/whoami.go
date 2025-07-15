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
		short = "Show the currently logged in user."
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
	_, err := api.AuthedClient()
	if err != nil {
		return err
	}

	username := "sanchitrk"

	fmt.Printf("%s\n", username)
	// user, err := client.Users.GetUser()
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%s\n", user.Username)
	//
	return nil
}
