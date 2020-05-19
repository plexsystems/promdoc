package rendering

import promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"

// RenderCSV renders CSV
func RenderCSV(ruleGroups []promv1.RuleGroup) string {
	var currentGroup string
	document := "Name,Summary,Message,Severity,Runbook\n"
	for _, ruleGroup := range ruleGroups {
		if currentGroup != ruleGroup.Name {
			currentGroup = ruleGroup.Name
			document += "\n## " + ruleGroup.Name + "\n"
		}

		for _, rule := range ruleGroup.Rules {
			if rule.Alert == "" {
				continue
			}

			document += rule.Alert + ";" + rule.Annotations["summary"] + ";" + trimSpaceNewlineInString(
				rule.Annotations["message"]) + ";" + rule.Labels["severity"] +
				";" + rule.Annotations["runbook_url"] + "\n"
		}
	}

	return document
}
