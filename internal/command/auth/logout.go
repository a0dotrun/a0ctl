package auth

import (
	"fmt"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/settings"

	"github.com/spf13/cobra"

	_ "embed"
)

func newLogout() *cobra.Command {
	const (
		use   = "logout"
		short = "Log out currently logged in user."
	)
	cmd := &cobra.Command{
		Use:               use,
		Short:             short,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cli.NoFilesArg,
		RunE:              logout,
		PersistentPreRunE: checkEnvAuth,
	}
	return cmd
}

func logout(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	settings, err := settings.ReadSettings()
	if err != nil {
		return fmt.Errorf("could not retrieve local config: %w", err)
	}

	if token := settings.GetToken(); len(token) == 0 {
		fmt.Println("No user logged in.")
		return nil
	}

	// if err := invalidateSessionsIfRequested(); err != nil {
	// 	return err
	// }

	settings.SetToken("")
	settings.SetUsername("")
	fmt.Println("Logged out.")

	return nil
}

// func invalidateSessionsIfRequested() error {
// 	if !flags.All() {
// 		return nil
// 	}

// 	client, err := api.AuthedClient()
// 	if err != nil {
// 		return err
// 	}

// 	from, err := client.Tokens.Invalidate()
// 	if err != nil {
// 		return err
// 	}

// 	formatted := time.Unix(from, 0).UTC().Format(time.DateTime)
// 	fmt.Printf("Invalidated all sessions started before %s UTC.\n", formatted)
// 	return nil
// }
