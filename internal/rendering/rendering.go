package rendering

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/ghodss/yaml"
)

var (
	errInvalidOutput = errors.New("output format not supported")
)

type TypeMeta struct {
	Kind string `json:"kind"`
}

type PrometheusRule struct {
	Spec PrometheusRuleSpec `json:"spec"`
}

type PrometheusRuleSpec struct {
	Groups []RuleGroup `json:"groups"`
}

type RuleGroup struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Record      string            `json:"record"`
	Alert       string            `json:"alert"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// Render renders rule groups in a specific output format
func Render(path string, outputType string) (string, error) {
	ruleGroups, err := getRuleGroups(path)
	if err != nil {
		return "", fmt.Errorf("get rule groups: %w", err)
	}

	switch outputType {
	case ".md":
		return RenderMarkdown(ruleGroups), nil
	case ".csv":
		return RenderCSV(ruleGroups), nil
	default:
		return "", errInvalidOutput
	}
}

func getRuleGroups(path string) ([]RuleGroup, error) {
	var ruleGroups []RuleGroup

	files, err := getYamlFiles(path)
	if err != nil {
		return nil, fmt.Errorf("get yaml files: %w", err)
	}

	for _, file := range files {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		var typeMeta TypeMeta
		if err := yaml.Unmarshal(fileContent, &typeMeta); err != nil {
			continue
		}

		if typeMeta.Kind != "PrometheusRule" {
			continue
		}

		var prometheusRule PrometheusRule
		if err := yaml.Unmarshal(fileContent, &prometheusRule); err != nil {
			continue
		}

		for _, group := range prometheusRule.Spec.Groups {
			var hasAlert bool
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

func trimText(text string) string {
	newline := regexp.MustCompile(`\r?\n|\| `)
	text = newline.ReplaceAllString(text, " ")

	prom := regexp.MustCompile(` \| `)
	text = prom.ReplaceAllString(text, " \\| ")

	return text
}
