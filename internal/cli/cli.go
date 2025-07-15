// Package cli provides the instance entry point for the root command.
package cli

import "github.com/spf13/cobra"

func NoFilesArg(
	cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}

// import (
// 	"context"
// 	"fmt"
// 	"os"
//
// 	"github.com/a0dotrun/a0ctl/internal/command/root"
// )
//
// func Run(ctx context.Context, args ...string) int {
// 	cmd := root.New()
// 	cmd.SetContext(ctx)
//
// 	cmd.SetArgs(args)
//
// 	if err := cmd.Execute(); err != nil {
// 		_, err := fmt.Fprintln(os.Stderr, err)
// 		if err != nil {
// 			return 0
// 		}
// 		return 1
// 	}
//
// 	return 0
// }
