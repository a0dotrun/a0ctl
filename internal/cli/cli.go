package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/a0dotrun/a0ctl/internal/command/root"
)

func Run(ctx context.Context, args ...string) int {

	cmd := root.New()
	cmd.SetContext(ctx)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil { // Should it be err == nil return 0
			return 0
		}
		return 1
	}

	return 0
}
