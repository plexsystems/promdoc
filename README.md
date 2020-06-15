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

`GO111MODULE=on go get github.com/plexsystems/promdoc`

## Usage

Promdoc will generate the output in the format that matches the output file.

For example, to generate markdown, run the following command in the root folder where you want `promdoc` to search for rules.

```console
$ promdoc generate alerts.md
```

To generate the output in `CSV`, you can run the same command, but rather than `.md`, use `.csv`

```console
$ promdoc generate alerts.csv
```

Supported output formats:

- Markdown (.md)
- CSV (.csv)

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

Running the generate command with a `.md` extension:

```console
$ promdoc generate alerts.md
```

Generates the following markdown:

# Alerts

## AlertManager

|Name|Summary|Severity|Runbook|
|---|---|---|---|
|AlertmanagerMembersInconsistent|Alertmanager has not found all other members of the cluster|critical|[https://runbook.com/alertmanagermembersinconsistent](https://runbook.com/alertmanagermembersinconsistent)|

## ControlPlane

|Name|Summary|Severity|Runbook|
|---|---|---|---|
|KubeControllerManagerDown|KubeControllerManager has disappeared from Prometheus target discovery|critical|[https://runbook.com/kubecontrollermanagerdown](https://runbook.com/kubecontrollermanagerdown)|
