// Package config provides commands to manage CLI configuration.
package config

import (
	"github.com/a0dotrun/a0ctl/internal/command/config/setconfig"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	const (
		short = "Manage your CLI configuration"
	)

	cmd := &cobra.Command{
		Use:   "config",
		Short: short,
	}

	cmd.AddCommand(
		setconfig.NewConfig(),
	)

	return cmd
}
