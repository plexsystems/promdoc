package generate

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

type typeMeta struct {
	Kind string `json:"kind"`
}

type prometheusRule struct {
	Spec prometheusRuleSpec `json:"spec"`
}

type prometheusRuleSpec struct {
	Groups []ruleGroup `json:"groups"`
}

type ruleGroup struct {
	Name  string `json:"name"`
	Rules []rule `json:"rules"`
}

type rule struct {
	Alert       string            `json:"alert"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// Generate finds all rules at the given path and its
// subdirectories and generates documentation with the
// specified file extension, including the period.
func Generate(path string, output string) (string, error) {
	switch output {
	case ".md":
		return Markdown(path)
	case ".csv":
		return CSV(path)
	default:
		return "", errors.New("output format not supported")
	}
}

func getRuleGroups(path string) ([]ruleGroup, error) {
	files, err := getYamlFiles(path)
	if err != nil {
		return nil, fmt.Errorf("get yaml files: %w", err)
	}

	var alertGroups []ruleGroup
	for _, file := range files {
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		var typeMeta typeMeta
		if err := yaml.Unmarshal(fileBytes, &typeMeta); err != nil {
			continue
		}
		if typeMeta.Kind != "PrometheusRule" {
			continue
		}

		var prometheusRule prometheusRule
		if err := yaml.Unmarshal(fileBytes, &prometheusRule); err != nil {
			continue
		}

		for _, group := range prometheusRule.Spec.Groups {
			var alertRules []rule
			for _, rule := range group.Rules {
				if rule.Alert == "" {
					continue
				}

				alertRules = append(alertRules, rule)
			}

			if len(alertRules) == 0 {
				continue
			}

			alertGroup := ruleGroup{
				Name:  group.Name,
				Rules: alertRules,
			}
			alertGroups = append(alertGroups, alertGroup)
		}
	}

	sort.Slice(alertGroups, func(i int, j int) bool {
		return alertGroups[i].Name < alertGroups[j].Name
	})

	return alertGroups, nil
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
