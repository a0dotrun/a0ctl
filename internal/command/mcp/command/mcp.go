package mcp

import (
	"github.com/a0dotrun/a0ctl/internal/command/mcp/builder"
	"github.com/a0dotrun/a0ctl/internal/command/mcp/initialize"
	"github.com/a0dotrun/a0ctl/internal/command/mcp/mcprun"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	const (
		short = "MCP (Model Context Protocol) commands"
		long  = "Commands for working with MCP applications including build, run, and initialization."
	)

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(
		mcpinitialize.New(),
		mcpbuilder.New(),
		mcprun.New(),
	)

	return cmd
}
