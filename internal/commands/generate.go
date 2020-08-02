package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/plexsystems/promdoc/internal/rendering"
)

// NewGenerateCommand creates a new generate command
func NewGenerateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "generate",
		Short: "Generate documentation from a given folder",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("out", cmd.Flags().Lookup("out")); err != nil {
				return fmt.Errorf("bind out flag: %w", err)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runGenerateCommand(); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("out", "o", "alerts.md",
		"file name or path for the output-file")

	return &cmd
}

func runGenerateCommand() error {
	out := viper.GetString("out")
	outputDir, outputFile := path.Split(out)
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}
	if outputDir == "" {
		outputDir = workingDir
	}
	if outputFile == "" {
		outputFile = "alerts.md"
	}

	fileExtension := path.Ext(outputFile)
	output, err := rendering.Render(workingDir, fileExtension)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	outputPath := path.Join(outputDir, outputFile)
	err = ioutil.WriteFile(outputPath, []byte(output), os.ModePerm)
	if err != nil {
		return fmt.Errorf("write document: %w", err)
	}

	return nil
}
