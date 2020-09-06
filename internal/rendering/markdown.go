package rendering

import (
	"fmt"
	"strings"
)

// RenderMarkdown renders Markdown.
func RenderMarkdown(ruleGroups []RuleGroup) string {
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

			document += "|Name|Summary|Description|Severity|Runbook|"
			document += "\n"

			document += "|---|---|---|---|---|"
			document += "\n"
		}

		for _, rule := range ruleGroup.Rules {
			if rule.Alert == "" {
				continue
			}

			var description string
			if val, ok := rule.Annotations["description"]; ok {
				description = trimText(val)
			} else if val, ok := rule.Annotations["message"]; ok {
				description = trimText(val)
			}

			summary := rule.Annotations["summary"]
			severity := rule.Labels["severity"]
			runbookURL := rule.Annotations["runbook_url"]
			if runbookURL != "" {
				runbookURL = fmt.Sprintf("[%s](%s)", runbookURL, runbookURL)
			}

			document += fmt.Sprintf("|%s|%s|%s|%s|%s|", rule.Alert, summary, description, severity, runbookURL)
			document += "\n"
		}
	}

	return document
}
