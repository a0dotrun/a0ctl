// Package setconfig helps with setting CLI configuration options.
package setconfig

import (
	"github.com/spf13/cobra"
)

func NewConfig() *cobra.Command {
	const (
		short = "Set a configuration value"
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: short,
	}

	cmd.AddCommand(
		newSetToken(),
	)

	return cmd
}
