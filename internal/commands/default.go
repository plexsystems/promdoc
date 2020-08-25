package commands

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

// NewDefaultCommand creates a new default command
func NewDefaultCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     path.Base(os.Args[0]),
		Short:   "Create documentation from Prometheus rules",
		Version: "0.5.0",
	}

	cmd.AddCommand(NewGenerateCommand())

	return &cmd
}
