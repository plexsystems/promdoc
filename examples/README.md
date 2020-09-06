# Examples

An example directory when working with `promdoc`.

## Alerting documentation

Each of the output type folders (e.g. CSV, Markdown) contain an `expected` file that is the result of running the `generate` command on this directory.

For example, to generate the output for Markdown, run the following command at the root of the repository:

```shell
promdoc generate examples --out examples/markdown/expected.md
```
