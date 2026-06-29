# go-ruby-yaml documentation

**Ruby's Psych-compatible YAML emitter and loader in pure Go — MRI-compatible, no cgo.**

`go-ruby-yaml/yaml` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's Psych YAML,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-yaml/yaml`.

It was **extracted from rbgo's prelude/internals into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: emitter + loader complete — Psych byte-exact"
    Faithful port of Ruby's Psych: **`Dump`** emitter and **`Load`** / **`SafeLoad`** loader, with **anchors and aliases**, the **`!ruby/object:`** Ruby tags, **`Symbol` / `Time` / `Range` / bignum** scalars, and both **block and flow** styles. Validated by a **differential oracle** against the system `ruby` / Psych — emitted YAML and loaded values compared byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
s, _ := yaml.Dump(map[string]any{"a": 1, "b": []any{2, 3}})
// "---\na: 1\nb:\n- 2\n- 3\n"

v, _ := yaml.Load("--- &a [1, 2]\n*a\n")   // shared via alias
v, _  = yaml.SafeLoad("--- !ruby/object:Foo {}")  // refused: error
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`yaml`](https://github.com/go-ruby-yaml/yaml) | the library — Ruby's Psych YAML in pure Go |
| [`docs`](https://github.com/go-ruby-yaml/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-yaml.github.io`](https://github.com/go-ruby-yaml/go-ruby-yaml.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-yaml/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-yaml/yaml](https://github.com/go-ruby-yaml/yaml).
