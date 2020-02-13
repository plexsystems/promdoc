package commands

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

// NewDefaultCommand creates a new default command
func NewDefaultCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   path.Base(os.Args[0]),
		Short: "promdoc",
		Long:  "A cli tool to create documentation from Prometheus rules",
	}

	cmd.AddCommand(NewGenerateCommand())

	return &cmd
}
