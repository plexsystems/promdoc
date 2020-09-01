package rendering

import (
	"fmt"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
)

// RenderCSV renders CSV
func RenderCSV(ruleGroups []promv1.RuleGroup) string {
	document := "Name,RuleGroup,Summary,Description,Severity,Runbook\n"
	for _, ruleGroup := range ruleGroups {
		for _, rule := range ruleGroup.Rules {
			if rule.Alert == "" {
				continue
			}

			var summary string
			if val, ok := rule.Annotations["summary"]; ok {
				summary = val
			} else if val, ok := rule.Annotations["message"]; ok {
				summary = replacePromQLInString(trimSpaceNewlineInString(val))
			}
			description := trimSpaceNewlineInString(rule.Annotations["description"])
			severity := rule.Labels["severity"]
			runbookURL := rule.Annotations["runbook_url"]

			document += fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", rule.Alert, ruleGroup.Name, summary, description, severity, runbookURL)
		}
	}

	return document
}
