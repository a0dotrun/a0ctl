// Package auth provides the authentication command ande helpers for the CLI.
package auth

import (
	"github.com/spf13/cobra"
)

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

	cmd.AddCommand(newWhoAMI(), newLogin())

	return cmd
}
