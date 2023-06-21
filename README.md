# Go Declarative Testing - Core ![go test workflow](https://github.com/jaypipes/gdt-core/actions/workflows/gate-tests.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/jaypipes/gdt-core.svg)](https://pkg.go.dev/github.com/jaypipes/gdt-core)

[`gdt`][gdt] is a testing library that allows test authors to cleanly describe tests
in a YAML file. `gdt` reads YAML files that describe a test's assertions and
then builds a set of Go structures that the standard Go
[`testing`](https://golang.org/pkg/testing/) package can execute.

[gdt]: https://github.com/jaypipes/gdt

This `gdt-core` repository is the core Go library for `gdt` that
contains the base types and plugin system.

## Contributing and acknowledgements

`gdt` was inspired by [Gabbi](https://github.com/cdent/gabbi), the excellent
Python declarative testing framework. `gdt` tries to bring the same clear,
concise test definitions to the world of Go functional testing.

Contributions to `gdt-core` are welcomed! Feel free to open a Github issue or
submit a pull request.
