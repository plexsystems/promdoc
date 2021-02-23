package generate

import (
	"fmt"
	"strings"
)

// Markdown finds all rules at the given path and its
// subdirectories and generates Markdown documentation
// from the found rules.
func Markdown(path string, input string) (string, error) {
	ruleGroups, err := getRuleGroups(path, input)
	if err != nil {
		return "", fmt.Errorf("get rule groups: %w", err)
	}

	document := "# Alerts"
	document += "\n\n"

	document += "## Rule Groups"
	document += "\n\n"
	var printingGroup string
	for _, ruleGroup := range ruleGroups {
		if printingGroup == ruleGroup.Name {
			continue
		}

		document += "* [" + ruleGroup.Name + "](#" + strings.ToLower(strings.ReplaceAll(ruleGroup.Name, " ", "-")) + ")"
		document += "\n"

		printingGroup = ruleGroup.Name
	}

	var currentGroup string
	for _, ruleGroup := range ruleGroups {
		if currentGroup != ruleGroup.Name {
			currentGroup = ruleGroup.Name
			document += "\n"
			document += "## " + ruleGroup.Name
			document += "\n\n"

			document += "|Name|Summary|Description|Severity|Expr|Runbook"
			document += "\n"

			document += "|---|---|---|---|---|---|"
			document += "\n"
		}

		for _, rule := range ruleGroup.Rules {
			var description string
			if val, ok := rule.Annotations["description"]; ok {
				description = trimText(val)
			} else if val, ok := rule.Annotations["message"]; ok {
				description = trimText(val)
			}

			expr := trimText(rule.Expr)
			// fmt.Printf("Expr: %s\n", expr)

			summary := rule.Annotations["summary"]
			severity := rule.Labels["severity"]
			runbookURL := rule.Annotations["runbook_url"]
			if runbookURL != "" {
				runbookURL = fmt.Sprintf("[%s](%s)", runbookURL, runbookURL)
			}

			document += fmt.Sprintf("|%s|%s|%s|%s|%s|%s", rule.Alert, summary, description, severity, expr, runbookURL)
			document += "\n"
		}
	}

	return document, nil
}
