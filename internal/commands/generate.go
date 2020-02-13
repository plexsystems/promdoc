package commands

import (
	"fmt"
	"io/ioutil"
	"os"
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
			return runGenerateCommand(args[0])
		},
	}

	return &cmd
}

func runGenerateCommand(output string) error {
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	files, err := getYamlFiles(path)
	if err != nil {
		return fmt.Errorf("get yaml files: %w", err)
	}

	prometheusRules, err := getPrometheusRules(files)
	if err != nil {
		return fmt.Errorf("get prometheus rules: %w", err)
	}

	sort.Slice(prometheusRules, func(i int, j int) bool {
		return prometheusRules[i].Spec.Groups[0].Name < prometheusRules[j].Spec.Groups[0].Name
	})

	var currentGroup string
	doc := "# Alerts"
	for _, prometheusRule := range prometheusRules {
		if currentGroup != prometheusRule.Spec.Groups[0].Name {
			currentGroup = prometheusRule.Spec.Groups[0].Name
			doc += "\n## " + prometheusRule.Spec.Groups[0].Name + "\n"
			doc += "|Name|Message|Severity|\n"
			doc += "|---|---|---|\n"
		}

		for _, rule := range prometheusRule.Spec.Groups[0].Rules {
			if rule.Alert == "" {
				continue
			}

			doc += "|" + rule.Alert + "|" + rule.Annotations["message"] + "|" + rule.Labels["severity"] + "|\n"
		}
	}

	err = ioutil.WriteFile(path+"/"+output, []byte(doc), os.ModePerm)
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

func getPrometheusRules(files []string) ([]v1.PrometheusRule, error) {
	var prometheusRules []v1.PrometheusRule

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

		prometheusRules = append(prometheusRules, prometheusRule)
	}

	return prometheusRules, nil
}
