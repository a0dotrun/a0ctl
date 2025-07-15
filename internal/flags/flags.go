// Package flags provides helpful flags for the CLI commands.
package flags

import (
	"github.com/spf13/cobra"
)

var resetConfig bool

func AddResetConfigFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&resetConfig, "reset-config", false, "")
	err := cmd.PersistentFlags().MarkHidden("reset-config")
	if err != nil {
		return
	}
}

func ResetConfig() bool {
	return resetConfig
}
