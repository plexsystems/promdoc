# promdoc

![logo](promdoc.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/plexsystems/promdoc)](https://goreportcard.com/report/github.com/plexsystems/promdoc)
[![GitHub release](https://img.shields.io/github/release/plexsystems/promdoc.svg)](https://github.com/plexsystems/promdoc/releases)

`promdoc` automatically generates documentation from your [PrometheusRules](https://github.com/coreos/prometheus-operator/blob/master/Documentation/design.md#prometheusrule).

## Installation

`GO111MODULE=on go get github.com/plexsystems/promdoc`

Binaries are also provided on the [releases](https://github.com/plexsystems/promdoc/releases) page.

## Usage

`promdoc` will generate the output in the format that matches the output file.

For example, to generate markdown, run the following command in the root folder where you want `promdoc` to search for rules.

```shell
$ promdoc generate
```

Optionally, you can specify a directory to generate alerts for. This will look at the specified directory and its subdirectories:

```shell
$ promdoc generate alertsdirectory
```

To generate the output in `CSV`, include the `.csv` extension in the output:

```shell
$ promdoc generate --out alerts.csv
```

**Supported output formats:**

- Markdown (.md)
- CSV (.csv)
