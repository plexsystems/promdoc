package generate

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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

func splitYAML(resources []byte) ([][]byte, error) {
	dec := yaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("decode: %w", err)
		}

		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshal: %w", err)
		}

		res = append(res, valueBytes)
	}

	return res, nil
}

func getMixinRuleGroups(files []string) ([]ruleGroup, error) {
	var alertGroups []ruleGroup

	for _, file := range files {
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		var byteSlices [][]byte
		byteSlices, err = splitYAML([]byte(fileBytes))
		if err != nil {
			return nil, fmt.Errorf("splitting yaml: %w", err)
		}

		for _, byteSlice := range byteSlices {
			var groups prometheusRuleSpec
			if err := yaml.Unmarshal(byteSlice, &groups); err != nil {
				continue
			}

			for _, group := range groups.Groups {
				alertGroup, err := extractGroupAlerts(group)
				if err != nil {
					return nil, fmt.Errorf("extract group alerts: %w", err)
				}

				alertGroups = append(alertGroups, *alertGroup)
			}
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

		var byteSlices [][]byte
		byteSlices, err = splitYAML([]byte(fileBytes))
		if err != nil {
			return nil, fmt.Errorf("splitting yaml: %w", err)
		}

		for _, byteSlice := range byteSlices {
			var typeMeta typeMeta
			if err := yaml.Unmarshal(byteSlice, &typeMeta); err != nil {
				continue
			}
			if typeMeta.Kind != "PrometheusRule" {
				continue
			}

			var prometheusRule prometheusRule
			if err := yaml.Unmarshal(byteSlice, &prometheusRule); err != nil {
				continue
			}

			for _, group := range prometheusRule.Spec.Groups {
				alertGroup, err := extractGroupAlerts(group)
				if err != nil {
					return nil, fmt.Errorf("extract group alerts: %w", err)
				}

				if alertGroup != nil {
					alertGroups = append(alertGroups, *alertGroup)
				}
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
