package rendering

import (
	"fmt"
	"strings"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
)

// RenderMarkdown renders Markdown
func RenderMarkdown(ruleGroups []promv1.RuleGroup) string {
	document := "# Alerts"
	document += "\n\n"

	document += "## Rule Groups"
	document += "\n\n"
	var printingGroup string
	for _, ruleGroup := range ruleGroups {
		if printingGroup == ruleGroup.Name {
			continue
		}

		document += "* [" + ruleGroup.Name + "](#" + strings.ReplaceAll(ruleGroup.Name, " ", "-") + ")"
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

			var summary string
			if val, ok := rule.Annotations["summary"]; ok {
				summary = val
			} else if val, ok := rule.Annotations["description"]; ok {
				summary = val
			} else if val, ok := rule.Annotations["message"]; ok {
				summary = replacePromQLInString(trimSpaceNewlineInString(val))
			}
			description := trimSpaceNewlineInString(rule.Annotations["description"])
			severity := rule.Labels["severity"]
			runbookURL := rule.Annotations["runbook_url"]
			if runbookURL != "" {
				runbookURL = fmt.Sprintf("[%s](%s)", runbookURL, runbookURL)
			}

			document += fmt.Sprintf("|%s|%s|%s|%s|%s", rule.Alert, summary, description, severity, runbookURL)
			document += "\n"
		}
	}

	return document
}
