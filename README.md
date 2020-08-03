# promdoc

![logo](promdoc.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/plexsystems/promdoc)](https://goreportcard.com/report/github.com/plexsystems/promdoc)
[![GitHub release](https://img.shields.io/github/release/plexsystems/promdoc.svg)](https://github.com/plexsystems/promdoc/releases)

`promdoc` automatically generates documentation from your [PrometheusRules](https://github.com/coreos/prometheus-operator/blob/master/Documentation/design.md#prometheusrule).

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

**Supported output formats:**

- Markdown (.md)
- CSV (.csv)
