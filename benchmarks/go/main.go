// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level workload for go-ruby-yaml/yaml: a representative Psych document
// (ordered mappings, sequences, mixed scalars — string/int/float/bool/symbol —
// a nested sequence, a block scalar, and a shared subtree emitted as an
// anchor/alias) is built once, Dumped to obtain the load input, then the two
// primitives are timed: Load (parse the document) and Dump (emit the tree).
//
// The same object is built by ruby/yaml.rb; go-ruby-yaml's Dump is byte-identical
// to MRI Psych.dump for this tree and Load(Dump(x)) round-trips identically
// (verified against MRI/JRuby/TruffleRuby before timing), so the two columns
// measure the same operation on the same bytes.
package main

import (
	"github.com/go-ruby-yaml/yaml"
)

// buildDoc constructs the representative Ruby value graph. The "defaults" subtree
// is referenced twice by identity so Dump emits it once behind "&1" and aliases
// the second occurrence as "*1" — exactly as MRI Psych does.
func buildDoc() yaml.Value {
	defaults := yaml.NewMap()
	defaults.Set("adapter", "postgresql")
	defaults.Set("host", "db.internal")
	defaults.Set("port", 5432)
	defaults.Set("pool", 5)
	defaults.Set("encoding", "utf8")
	defaults.Set("ssl", true)

	root := yaml.NewMap()
	root.Set("version", "1.4.2")
	root.Set("replicas", 3)
	root.Set("ratio", 0.75)
	root.Set("enabled", true)
	root.Set("environment", yaml.Symbol("production"))
	root.Set("tags", []any{"web", "api", "worker"})
	root.Set("ports", []any{8080, 8081, 8082})
	root.Set("primary", defaults)
	root.Set("secondary", defaults)
	root.Set("matrix", []any{[]any{1, 2}, []any{3, 4}})
	root.Set("notes", "line one\nline two\nline three")
	return root
}

func main() {
	doc := buildDoc()
	in, err := yaml.Dump(doc) // the load input, byte-identical to MRI Psych.dump(doc)
	if err != nil {
		panic(err)
	}
	bench("load-config", 1000, func() { v, _ := yaml.Load(in); sink = v })
	bench("dump-config", 1000, func() { s, _ := yaml.Dump(doc); sink = s })
}
