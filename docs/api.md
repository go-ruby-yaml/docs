# Usage & API

The public API lives at the module root (`github.com/go-ruby-yaml/yaml`). It is **Ruby-shaped but Go-idiomatic**: `Dump` / `Load` mirror Psych's `Psych.dump` / `Psych.load`, while the surface follows Go conventions — an explicit `error`, value types, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-yaml/yaml`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-yaml/yaml
```

## Worked example

```go
s, _ := yaml.Dump(map[string]any{"a": 1, "b": []any{2, 3}})
// "---\na: 1\nb:\n- 2\n- 3\n"

v, _ := yaml.Load("--- &a [1, 2]\n*a\n")   // shared via alias
v, _  = yaml.SafeLoad("--- !ruby/object:Foo {}")  // refused: error
```

## Shape

```go
// Dump emits Psych-compatible YAML for a Ruby value graph.
func Dump(v Value) (string, error)

// Load parses YAML into the full Ruby value graph (Psych.load).
func Load(s string) (Value, error)

// SafeLoad parses only the safe scalar/collection subset,
// refusing arbitrary !ruby/object: instantiation.
func SafeLoad(s string) (Value, error)
```

## Conformance

The loader is validated against the upstream
[yaml-test-suite](https://github.com/yaml/yaml-test-suite) on its
schema-independent **accept/reject** axis and passes **all 402/402 tests** —
**308/308 well-formed documents load** and **94/94 ill-formed inputs are
rejected** (0 panics, 0 known gaps). A ratchet test fails on any regression.

[![yaml-test-suite 402/402](https://img.shields.io/badge/yaml--test--suite-402%2F402-1a7f37)](https://github.com/yaml/yaml-test-suite)

Byte-exact correctness is additionally defined by reference Ruby: a **differential
oracle** runs a wide corpus through both the system `ruby` and this library and
compares the results **byte-for-byte** — not approximated from memory. The oracle
tests skip themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so
the cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-yaml/yaml` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
