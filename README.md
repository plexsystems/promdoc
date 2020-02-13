# promdoc

`promdoc` lets you automatically generate documentation from your [PrometheusRules](https://github.com/coreos/prometheus-operator/blob/master/Documentation/design.md#prometheusrule).

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

## Example

Given the following `PrometheusRule` definitions:

*kubecontrollermanagerdown-alert.yaml*
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    prometheus: k8s
    role: alert-rules
  name: kubecontrollermanagerdown-alerts
spec:
  groups:
  - name: ControlPlane
    rules:
    - alert: KubeControllerManagerDown
      annotations:
        message: KubeControllerManager has disappeared from Prometheus target discovery.
        runbook_url: https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubecontrollermanagerdown
      expr: |
        absent(up{job="kube-controller-manager"}) > 0
      for: 15m
      labels:
        severity: critical
```

*alertmanagermembersinconsistent-alerts.yaml*
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    role: alert-rules
  name: alertmanagermembersinconsistent-alerts
  namespace: monitoring
spec:
  groups:
  - name: AlertManager
    rules:
    - alert: AlertmanagerMembersInconsistent
      annotations:
        message: Alertmanager has not found all other members of the cluster.
      expr: |
        alertmanager_cluster_members{job="alertmanager-main",namespace="monitoring"}
          != on (service) GROUP_LEFT()
        count by (service) (alertmanager_cluster_members{job="alertmanager-main",namespace="monitoring"})
      for: 5m
      labels:
        severity: critical
```

The generated documentation would be

## Alerts

## AlertManager
|Name|Message|Severity|
|---|---|---|
|AlertmanagerMembersInconsistent|Alertmanager has not found all other members of the cluster.|critical|

## ControlPlane
|Name|Message|Severity|
|---|---|---|
|KubeControllerManagerDown|KubeControllerManager has disappeared from Prometheus target discovery.|critical|



