package rendering

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/ghodss/yaml"
)

var (
	errInvalidOutput = errors.New("output format not supported")
)

// Render renders rule groups in a specific output format
func Render(path string, outputType string) (string, error) {
	ruleGroups, err := getRuleGroups(path)
	if err != nil {
		return "", fmt.Errorf("get rule groups: %w", err)
	}

	switch outputType {
	case "md":
		return RenderMarkdown(ruleGroups), nil
	case "csv":
		return RenderCSV(ruleGroups), nil
	default:
		return "", errInvalidOutput
	}
}

func getRuleGroups(path string) ([]v1.RuleGroup, error) {
	var ruleGroups []v1.RuleGroup

	files, err := getYamlFiles(path)
	if err != nil {
		return nil, fmt.Errorf("get yaml files: %w", err)
	}

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

	sort.Slice(ruleGroups, func(i int, j int) bool {
		return ruleGroups[i].Name < ruleGroups[j].Name
	})

	return ruleGroups, nil
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

func trimSpaceNewlineInString(s string) string {
	re := regexp.MustCompile(`\r?\n|\| `)
	return re.ReplaceAllString(s, " ")
}
