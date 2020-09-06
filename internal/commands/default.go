package commands

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

// version is the version of the CLI, and it is set at build time.
var version string

// NewDefaultCommand creates a new default command.
func NewDefaultCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     path.Base(os.Args[0]),
		Short:   "Create documentation from Prometheus rules",
		Version: version,
	}

	cmd.AddCommand(NewGenerateCommand())

	return &cmd
}
