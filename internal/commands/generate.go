package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"github.com/plexsystems/promdoc/internal/rendering"
)

var (
	errInvalidOutput = errors.New("Output format not supported")
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
				return fmt.Errorf("Invalid argument: %s", err)
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

	files, err := getYamlFiles(workingDir)
	if err != nil {
		return fmt.Errorf("get yaml files: %w", err)
	}

	ruleGroups, err := getRuleGroups(files)
	if err != nil {
		return fmt.Errorf("get prometheus rules: %w", err)
	}

	sort.Slice(ruleGroups, func(i int, j int) bool {
		return ruleGroups[i].Name < ruleGroups[j].Name
	})

	document, err := getDocumentation(ruleGroups, outputType)
	if err != nil {
		return err
	}
	outputPath := path.Join(workingDir, outputFile)
	err = ioutil.WriteFile(outputPath, []byte(document), os.ModePerm)
	if err != nil {
		return fmt.Errorf("write document: %w", err)
	}

	return nil
}

func getYamlFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(currentFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk path: %w", err)
		}

		if fileInfo.IsDir() && fileInfo.Name() == ".git" {
			return filepath.SkipDir
		}

		if fileInfo.IsDir() {
			return nil
		}

		if filepath.Ext(currentFilePath) != ".yaml" && filepath.Ext(currentFilePath) != ".yml" {
			return nil
		}

		files = append(files, currentFilePath)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func getRuleGroups(files []string) ([]v1.RuleGroup, error) {
	var ruleGroups []v1.RuleGroup

	for _, file := range files {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		var prometheusRule v1.PrometheusRule
		err = yaml.Unmarshal(fileContent, &prometheusRule)
		if err != nil {
			continue
		}

		if prometheusRule.Kind != "PrometheusRule" {
			continue
		}

		for _, group := range prometheusRule.Spec.Groups {
			hasAlert := false

			for _, rule := range group.Rules {
				if rule.Alert != "" {
					hasAlert = true
					break
				}
			}

			if hasAlert {
				ruleGroups = append(ruleGroups, group)
			}
		}
	}

	return ruleGroups, nil
}

func getDocumentation(ruleGroups []v1.RuleGroup, outputType string) (string, error) {
	switch outputType {
	case "markdown":
		return rendering.RenderMarkdown(ruleGroups), nil
	case "csv":
		return rendering.RenderCSV(ruleGroups), nil
	default:
		return "", errInvalidOutput
	}

}
