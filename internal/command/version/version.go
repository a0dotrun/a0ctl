package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

// New initializes and returns a new version Command.
func New() *cobra.Command {
	const (
		short = "Show version information for the a0ctl command"
		long  = "Shows version information for the a0ctl command itself - including version number and builder date."
	)

	return &cobra.Command{
		Use:   "version",
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("a0ctl version %s\n", "0.0.1")
		},
	}
}
