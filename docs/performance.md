# Performance

`go-ruby-yaml/yaml` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `yaml`. This
page records a **comparative benchmark** of that module against the reference
Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## What is measured

The **same** Ruby script — `YAML.dump` + `YAML.load` round-trip of a representative config structure — is run under every runtime. `rbgo`'s
number reflects **this pure-Go library doing the work**; every other column is
that interpreter's own `yaml` stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI**
before timing.

- **Host:** Apple M4 Max, macOS (darwin/arm64). **Method:** best-of-5 wall time
  (best, not mean, to suppress scheduler noise); single-shot processes, no
  warm-up beyond the script's own loop.
- **Runtimes:** `ruby 4.0.5 +PRISM` (MRI, the oracle) and `ruby --yjit`;
  `jruby 10.1.0.0` (OpenJDK 25); `truffleruby 34.0.1` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`yaml.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-yaml) | 170 | 0.23× |
| MRI (ruby 4.0.5) | 750 | 1.00× |
| MRI + YJIT | 480 | 0.64× |
| JRuby 10.1.0.0 | 2460 | 3.28× |
| TruffleRuby 34.0.1 | 3370 | 4.49× |

rbgo runs on **go-ruby-yaml** and is **~4x faster than MRI** here (0.23x). MRI's YAML round-trip is comparatively heavy; the compiled pure-Go dump/load wins clearly. TruffleRuby pays heavy cold warm-up on this row (3370 ms).

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are real measured numbers from the
    2026-06-29 run — nothing is cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `yaml`
(Psych, wrapping C libyaml)?* The **same workload, same inputs, same iteration
counts** run through the Go library and through each reference runtime's stdlib;
emitted YAML and loaded structure were checked **byte-identical to MRI** (and to
JRuby / TruffleRuby) before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` (Psych 5.3.1 / libyaml 0.2.5)
  · MRI + YJIT · JRuby 10.1.0.0 (OpenJDK 25, Psych 5.3.1) · TruffleRuby 34.0.1
  (GraalVM CE Native, Psych 5.2.2).
- **Workload:** a representative Psych document — an ordered mapping of mixed
  scalars (string / integer / float / boolean / symbol), string and integer
  sequences, a nested sequence, a multi-line block scalar, and a **shared
  subtree** emitted once behind an anchor (`&1`) and aliased (`*1`).
  `load-config` parses it; `dump-config` emits it. Ruby's `YAML.load` runs with
  `aliases: true` (Psych 5 requires it for aliased documents).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### dump-config

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 3883.4 | 0.07× |
| MRI | 54223.0 | 1.00× |
| MRI + YJIT | 28934.0 | 0.53× |
| JRuby | 32377.5 | 0.60× |
| TruffleRuby | 145228.5 | 2.68× |

#### load-config

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 13834.4 | 0.23× |
| MRI | 60235.0 | 1.00× |
| MRI + YJIT | 36541.0 | 0.61× |
| JRuby | 37797.8 | 0.63× |
| TruffleRuby | 117752.3 | 1.95× |

The task's caveat was that Psych wraps **C libyaml**, so parity might be hard —
but the pure-Go library is in fact **~14× faster than MRI on dump** (0.07×) and
**~4.3× faster on load** (0.23×). The reason is that Psych's cost is *not* in
libyaml's C scanning/emitting: it is in building and walking the Ruby node tree
(a visitor allocating a Psych AST plus the target object graph, per call). The
lean pure-Go emitter writes bytes straight from the value graph and the pure-Go
loader builds the ordered map directly, so both primitives win outright. Every
runtime's output was verified byte-identical to MRI before timing — the anchor
(`&1`) / alias (`*1`), the `|-` block scalar, and the `:production` symbol all
round-trip exactly. TruffleRuby is a **cold-JIT outlier** on these short loops
(slower than plain MRI); read its column as order-of-magnitude, not steady state.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-yaml/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/yaml.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby, whose cold Graal JIT lands slower than
    MRI on these short loops. Every number here is a **real measured value** from
    the dated run above — nothing is fabricated, estimated, or cherry-picked. The
    go-ruby column is the pure-Go library; every other column is that
    interpreter's own Psych stdlib doing the equivalent work.
