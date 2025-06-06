package root

import (
	"log"
	"os"
	"path/filepath"

	"github.com/a0dotrun/a0ctl/internal/command/builder"
	"github.com/a0dotrun/a0ctl/internal/command/initialize"
	"github.com/a0dotrun/a0ctl/internal/command/run"
	"github.com/a0dotrun/a0ctl/internal/command/version"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	const (
		long  = "This is a0ctl - the a0.run command line interface."
		short = "The a0.run command line interface"
	)

	exePath, err := os.Executable()
	var exe string
	if err != nil {
		log.Printf("WARN: failed to find executable, error=%q", err)
		exe = "a0"
	} else {
		exe = filepath.Base(exePath)
	}

	root := &cobra.Command{
		Use: exe, Short: short, Long: long,
	}

	root.AddCommand(
		version.New(),
		initialize.New(),
		builder.New(),
		run.New(),
	)

	return root
}
