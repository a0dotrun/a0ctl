// Package cli provides the instance entry point for the root command.
package cli

import "github.com/spf13/cobra"

func NoFilesArg(
	_ *cobra.Command, _ []string, _ string,
) ([]string, cobra.ShellCompDirective) {
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}
