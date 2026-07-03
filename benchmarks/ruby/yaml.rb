# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference workload mirroring benchmarks/go/main.go: the same representative
# Psych document is built, dumped once to obtain the load input, then Load and
# Dump are timed. YAML.load is called with aliases: true because the document
# uses an anchor/alias (Psych 5 rejects aliases by default). The built object and
# the emitted YAML are byte-identical to go-ruby-yaml's (verified before timing).
require "yaml"
require_relative "_harness"

def build_doc
  defaults = {
    "adapter"  => "postgresql",
    "host"     => "db.internal",
    "port"     => 5432,
    "pool"     => 5,
    "encoding" => "utf8",
    "ssl"      => true,
  }
  {
    "version"     => "1.4.2",
    "replicas"    => 3,
    "ratio"       => 0.75,
    "enabled"     => true,
    "environment" => :production,
    "tags"        => ["web", "api", "worker"],
    "ports"       => [8080, 8081, 8082],
    "primary"     => defaults,
    "secondary"   => defaults,
    "matrix"      => [[1, 2], [3, 4]],
    "notes"       => "line one\nline two\nline three",
  }
end

doc = build_doc
input = YAML.dump(doc) # byte-identical to go-ruby-yaml Dump(doc)
bench("load-config", 1000) { YAML.load(input, aliases: true) }
bench("dump-config", 1000) { YAML.dump(doc) }
