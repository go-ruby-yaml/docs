# Roadmap

`go-ruby-yaml/yaml` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's psych yaml — the
deterministic, interpreter-independent slice extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Scanner & parser | A YAML 1.1 scanner and parser following Psych's event model — documents, mappings, sequences, scalars in plain / single / double / literal / folded styles, and the directive and document markers. | **Done** |
| Dump emitter | `Dump(v)` emits Psych-compatible YAML: block style by default with flow where Psych uses it, correct scalar quoting, and stable key ordering — matching reference Ruby's output byte-for-byte. | **Done** |
| Load & SafeLoad | `Load(s)` builds the full Ruby value graph; `SafeLoad` restricts to the safe scalar/collection subset, refusing arbitrary `!ruby/object:` instantiation exactly as Psych's safe loader does. | **Done** |
| Anchors, aliases & Ruby tags | Anchors `&a` and aliases `*a` for shared and cyclic structure, plus the `!ruby/object:`, `!ruby/symbol`, `!ruby/range` and related tags round-tripped the way Psych writes and reads them. | **Done** |
| Ruby scalar types | `Symbol`, `Time`, `Range` and arbitrary-precision bignum scalars emitted and parsed with Psych's tag and formatting conventions. | **Done** |
| Differential oracle & coverage | A wide value corpus dumped and loaded both here and by the system `ruby`/Psych, compared byte-for-byte; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
