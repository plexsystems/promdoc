package rendering

import (
	"fmt"
)

// RenderCSV renders CSV
func RenderCSV(ruleGroups []RuleGroup) string {
	document := "Name,RuleGroup,Summary,Description,Severity,Runbook\n"
	for _, ruleGroup := range ruleGroups {
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

			document += fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", rule.Alert, ruleGroup.Name, summary, description, severity, runbookURL)
		}
	}

	return document
}
