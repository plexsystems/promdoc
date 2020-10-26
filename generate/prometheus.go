package generate

import (
	"fmt"
	"io/ioutil"

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

func getMixinRuleGroups(files []string) ([]ruleGroup, error) {
	var alertGroups []ruleGroup

	for _, file := range files {
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		var groups prometheusRuleSpec
		if err := yaml.Unmarshal(fileBytes, &groups); err != nil {
			continue
		}

		for _, group := range groups.Groups {
			alertGroup, err := extractGroupAlerts(group)
			if err != nil {
				return nil, err
			}
			alertGroups = append(alertGroups, *alertGroup)
		}

	}

	return alertGroups, nil
}

func getKubernetesRuleGroups(files []string) ([]ruleGroup, error) {
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
			alertGroup, err := extractGroupAlerts(group)
			if err != nil {
				return nil, err
			}
			if alertGroup != nil {
				alertGroups = append(alertGroups, *alertGroup)
			}
		}

	}
	return alertGroups, nil
}

func extractGroupAlerts(group ruleGroup) (*ruleGroup, error) {
	var alertRules []rule
	for _, rule := range group.Rules {
		if rule.Alert == "" {
			return nil, nil
		}

		alertRules = append(alertRules, rule)
	}

	if len(alertRules) == 0 {
		return nil, nil
	}

	alertGroup := ruleGroup{
		Name:  group.Name,
		Rules: alertRules,
	}
	return &alertGroup, nil
}
