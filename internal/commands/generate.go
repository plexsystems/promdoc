package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
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

	document := getDocumentation(ruleGroups)
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

func getDocumentation(ruleGroups []v1.RuleGroup) string {
	document := "# Alerts"

	var currentGroup string
	for _, ruleGroup := range ruleGroups {
		if currentGroup != ruleGroup.Name {
			currentGroup = ruleGroup.Name
			document += "\n## " + ruleGroup.Name + "\n"
			document += "|Name|Message|Severity|\n"
			document += "|---|---|---|\n"
		}

		for _, rule := range ruleGroup.Rules {
			if rule.Alert == "" {
				continue
			}

			document += "|" + rule.Alert + "|" + rule.Annotations["message"] + "|" + rule.Labels["severity"] + "|\n"
		}
	}

	return document
}
