# promdoc

![logo](promdoc.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/plexsystems/promdoc)](https://goreportcard.com/report/github.com/plexsystems/promdoc)
[![GitHub release](https://img.shields.io/github/release/plexsystems/promdoc.svg)](https://github.com/plexsystems/promdoc/releases)

`promdoc` automatically generates documentation from your [PrometheusRules](https://github.com/coreos/prometheus-operator/blob/master/Documentation/design.md#prometheusrule).

```
NOTE: This project is currently a work in progress.
Feedback, feature requests, and contributions are welcome!
```

## Installation

### From source

`go get github.com/plexsystems/promdoc`

## Usage

Run the following command in the root of your project to create a markdown file named `alerts.md` that contains your Prometheus alerts.

```console
$ promdoc generate alerts.md
```

**NOTE:** The `summary` annotation on the `PrometheusRule` CRD is used for generating the summary on the document.

## Example

Given the following `PrometheusRule` definitions:

*controlplane-alerts.yaml*
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: controlplane-alerts
spec:
  groups:
  - name: ControlPlane
    rules:
    - alert: KubeControllerManagerDown
      annotations:
        summary: KubeControllerManager has disappeared from Prometheus target discovery
        runbook_url: https://runbook.com/kubecontrollermanagerdown
      expr: expression()
      labels:
        severity: critical
```

*alertmanager-alerts.yaml*
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: alertmanager-alerts
spec:
  groups:
  - name: AlertManager
    rules:
    - alert: AlertmanagerMembersInconsistent
      annotations:
        summary: Alertmanager has not found all other members of the cluster
        runbook_url: https://runbook.com/alertmanagermembersinconsistent
      expr: expression()
      labels:
        severity: critical
```

The generated documentation would be

# Alerts

## AlertManager

|Name|Summary|Severity|Runbook|
|---|---|---|---|
|AlertmanagerMembersInconsistent|Alertmanager has not found all other members of the cluster|critical|[https://runbook.com/alertmanagermembersinconsistent](https://runbook.com/alertmanagermembersinconsistent)|

## ControlPlane

|Name|Summary|Severity|Runbook|
|---|---|---|---|
|KubeControllerManagerDown|KubeControllerManager has disappeared from Prometheus target discovery|critical|[https://runbook.com/kubecontrollermanagerdown](https://runbook.com/kubecontrollermanagerdown)|
