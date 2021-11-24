package generate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// Generate finds all rules at the given path and its
// subdirectories and generates documentation with the
// specified file extension, including the period.
func Generate(path string, output string, input string) (string, error) {
	switch output {
	case ".md":
		return Markdown(path, input)
	case ".csv":
		return CSV(path, input)
	default:
		return "", errors.New("output format not supported")
	}
}

func getRuleGroups(path string, input string) ([]ruleGroup, error) {
	files, err := getYamlFiles(path)
	if err != nil {
		return nil, fmt.Errorf("get yaml files: %w", err)
	}

	var alertGroups []ruleGroup
	if input == "kubernetes" {
		alertGroups, err = getKubernetesRuleGroups(files)
		if err != nil {
			return nil, err
		}
	} else {
		alertGroups, err = getMixinRuleGroups(files)
		if err != nil {
			return nil, err
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
