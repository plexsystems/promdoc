package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/plexsystems/promdoc/internal/rendering"
)

// NewGenerateCommand creates a new generate command
func NewGenerateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "generate <directory>",
		Short: "Generate documentation from a given folder",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("out", cmd.Flags().Lookup("out")); err != nil {
				return fmt.Errorf("bind out flag: %w", err)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			if err := runGenerateCommand(path); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("out", "o", "alerts.md", "File name or directory for the alert documentation")

	return &cmd
}

func runGenerateCommand(path string) error {
	outputPath := viper.GetString("out")
	if filepath.Ext(outputPath) == "" {
		outputPath = filepath.Join(outputPath, "alerts.md")
	}

	output, err := rendering.Render(path, filepath.Ext(outputPath))
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, []byte(output), os.ModePerm); err != nil {
		return fmt.Errorf("write document: %w", err)
	}

	return nil
}
