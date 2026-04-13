# Project Structure

This structure is mandatory unless a documented architectural reason requires change.

```text
sneaky-core/
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── docs/
│   ├── 01_PRODUCT_SCOPE.md
│   ├── 02_ARCHITECTURE.md
│   ├── 03_PROJECT_STRUCTURE.md
│   ├── 04_PHASE_PLAN.md
│   ├── 05_PROTOCOL_MATRIX.md
│   ├── 06_CLI_SPEC.md
│   ├── 07_BINDING_PREPARATION.md
│   ├── PROTOCOL_MATRIX.md
│   ├── SINGBOX_COVERAGE_RESULTS.md
│   ├── prompts/
│   │   └── SINGBOX_COVERAGE_PHASE_PROMPT.md
│   └── AI_DEV_RULES.md
├── cmd/
│   └── sneakycli/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── manager.go
│   │   ├── state.go
│   │   └── errors.go
│   ├── config/
│   │   ├── detector.go
│   │   ├── parser.go
│   │   ├── validator.go
│   │   └── types.go
│   ├── cli/
│   │   ├── app.go
│   │   └── state.go
│   ├── adapter/
│   │   ├── adapter.go
│   │   ├── registry.go
│   │   └── capabilities.go
│   ├── adapters/
│   │   ├── singbox/
│   │       ├── adapter.go
│   │       └── validator.go
│   │   └── ssh/
│   │       ├── adapter.go
│   │       └── config.go
│   ├── runtime/
│   │   ├── session.go
│   │   ├── lifecycle.go
│   │   ├── context.go
│   │   └── process.go
│   ├── stats/
│   │   ├── stats.go
│   │   └── snapshot.go
│   └── logx/
│       ├── logger.go
│       └── entries.go
├── pkg/
│   └── sneaky/
│       └── sneaky.go
├── testdata/
│   ├── certs/
│   └── singbox/
│       ├── http/
│       ├── hysteria/
│       ├── hysteria2/
│       ├── shadowsocks/
│       ├── trojan/
│       ├── tuic/
│       ├── utilities/
│       ├── vless/
│       ├── vmess/
│       └── wireguard/
├── tests/
│   ├── config/
│   ├── adapters/
│   └── integration/
│       └── singbox_coverage_test.go
└── scripts/
    ├── build.sh
    └── test.sh
```

## Structure Rules

- `cmd/` is only for executable entrypoints
- `pkg/sneaky/` is the stable external package surface
- `internal/` holds implementation details
- protocol-specific code must stay inside `internal/adapters/`
- shared abstractions belong in `internal/adapter/`
- tests must reflect architecture, not random file placement

## Forbidden Structure Problems

Do not:
- put everything in one package
- mix CLI logic with adapter logic
- mix config detection with runtime control
- store experimental junk in root
- create duplicate config models in multiple packages
