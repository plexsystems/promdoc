package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/plexsystems/promdoc/internal/rendering"
)

// NewGenerateCommand creates a new generate command
func NewGenerateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "generate <output-dir>",
		Short: "Generate documentation from a given folder",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return fmt.Errorf("invalid argument: %w", err)
			}

			if err := runGenerateCommand(args[0], output); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			return nil
		},
	}
	cmd.Flags().String("output", "markdown", "Output format: markdown, csv")

	return &cmd
}

func runGenerateCommand(outputFile string, outputType string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	output, err := rendering.Render(workingDir, outputType)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	outputPath := path.Join(workingDir, outputFile)
	err = ioutil.WriteFile(outputPath, []byte(output), os.ModePerm)
	if err != nil {
		return fmt.Errorf("write document: %w", err)
	}

	return nil
}
