package root

import (
	"github.com/a0dotrun/a0ctl/internal/command/version"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
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

	root.AddCommand(version.New())

	return root
}
