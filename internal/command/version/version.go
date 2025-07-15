// Package version provides the command to show the version of a0ctl.
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Show version information for the a0ctl CLI."
		long  = "Shows version information for the a0ctl CLI."
	)

	return &cobra.Command{
		Use:   "version",
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			_ = ctx // Access the context here
			fmt.Printf("a0ctl version %s\n", "0.0.1")
		},
	}
}
