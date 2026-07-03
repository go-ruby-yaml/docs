<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-yaml` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-yaml` library**
against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby, TruffleRuby). It
measures the **library primitive** through its Go API, isolated from the rbgo
interpreter, so the numbers answer: *is the pure-Go implementation as fast as the
reference runtime's own Psych?*

## Layout

- `go/`          — self-contained Go driver; `go.mod` pins the published library.
- `ruby/yaml.rb` — the equivalent workload; `ruby/_harness.rb` is the shared timer.
- `run.sh`       — runs every available runtime and prints one Markdown table per
  sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Workload

A representative Psych document: an ordered mapping of mixed scalars
(string / integer / float / boolean / symbol), string and integer sequences, a
nested sequence, a multi-line block scalar, and a **shared subtree** referenced
twice so it is emitted once behind an anchor (`&1`) and aliased (`*1`). Two
primitives are timed:

- **`load-config`** — parse the document into the Ruby value graph
  (`YAML.load` / `yaml.Load`; Ruby uses `aliases: true` as Psych 5 requires).
- **`dump-config`** — emit the value graph back to a Psych document
  (`YAML.dump` / `yaml.Dump`).

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build **identical inputs** and
their emitted YAML is **byte-identical to MRI Psych** (also verified against
JRuby and TruffleRuby) before timing — `go-ruby-yaml`'s `Dump` matches
`Psych.dump` byte-for-byte and `Load(Dump(x))` round-trips identically. Results
are published, dated, in `../docs/performance.md`.
