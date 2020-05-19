package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

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
			if err := runGenerateCommand(args[0]); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			return nil
		},
	}

	return &cmd
}

func runGenerateCommand(outputFile string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	fileTokens := strings.Split(outputFile, ".")
	if len(fileTokens) == 0 {
		return fmt.Errorf("get file extension: %w", err)
	}

	fileExtension := fileTokens[len(fileTokens)-1]
	output, err := rendering.Render(workingDir, fileExtension)
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
