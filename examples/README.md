# Examples

An example directory when working with `Promdoc`.

## Alerting documentation

Each of the output type folders (e.g. csv, markdown) contain an `expected` file that is the result of running the `generate` command on [rule.yaml](rule.yaml).

For example, to generate the output for _Markdown_, run the following command at the root of the repository:

```shell
promdoc generate examples --out examples/markdown/expected.md
```
