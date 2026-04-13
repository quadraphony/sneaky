# Phase Plan

This project must be built in controlled phases.

A phase cannot be skipped.

---

## Phase 0 — Foundation and Repo Setup

### Objective
Create a clean, strict, build-ready repo foundation.

### Deliverables
- repository initialized
- folder structure created
- Go module initialized
- Makefile added
- core docs added
- lint and test scripts added
- initial empty but valid package layout created

### Acceptance Criteria
- repo builds at skeleton level
- no broken imports
- docs match actual structure

### Do Not Do In This Phase
- no protocol logic
- no fake handlers
- no mobile binding work

---

## Phase 1 — Core Contracts

### Objective
Define the core interfaces and lifecycle model.

### Deliverables
- public core manager contract
- adapter interface
- runtime state model
- structured error model
- start and stop lifecycle rules
- adapter registry contract

### Acceptance Criteria
- code compiles
- lifecycle rules are explicit
- contract tests pass

### Do Not Do In This Phase
- no actual sing-box execution yet
- no multi-protocol claims yet

---

## Phase 2 — Config Detection and Validation

### Objective
Build deterministic config ingestion and detection.

### Deliverables
- raw config ingestion
- format detection
- typed config metadata
- validation errors
- ambiguity handling rules
- detection unit tests

### Acceptance Criteria
- invalid input fails clearly
- ambiguous detection fails clearly
- tests cover expected config classes

### Do Not Do In This Phase
- no silent fallback logic
- no protocol expansion beyond defined detection scope

---

## Phase 3 — Sing-box Adapter Foundation

### Objective
Integrate the first real adapter path.

### Deliverables
- sing-box adapter package
- adapter validation path
- startup and shutdown flow
- runtime lifecycle hooks
- adapter registration
- integration tests where possible

### Acceptance Criteria
- manager can resolve and run sing-box adapter
- stop path is clean
- errors are surfaced clearly

### Do Not Do In This Phase
- no second adapter yet
- no mobile bindings yet

---

## Phase 4 — Logging and Stats

### Objective
Expose observability consistently.

### Deliverables
- structured logging pipeline
- stats snapshot model
- runtime status exposure
- CLI-readable output
- test coverage for non-runtime pieces

### Acceptance Criteria
- logs are structured and useful
- stats access is stable
- failures remain visible

---

## Phase 5 — CLI Tooling

### Objective
Create a real CLI for testing and development workflows.

### Deliverables
- `start` command
- `stop` command if applicable to process model
- `status` command
- `validate` command
- clear CLI exit codes
- CLI usage documentation

### Acceptance Criteria
- CLI can validate configs
- CLI can start supported adapter path
- CLI output is readable and actionable

---

## Phase 6 — Second Adapter Introduction

### Objective
Add one additional adapter only after the architecture proves stable.

### Candidate Options
Choose one only after confirming implementation path:
- OpenVPN adapter
- SSH-based adapter
- WireGuard adapter

### Rule
Do not start more than one new adapter in the same phase.

### Acceptance Criteria
- second adapter integrates through existing contracts
- no core rewrite is needed
- config detection extends safely

---

## Phase 7 — Binding Preparation

### Objective
Prepare the core for later Android and iOS integration.

### Deliverables
- stable exported package surface
- binding-safe API review
- reduced unnecessary API exposure
- lifecycle suitability for bindings
- documentation for bridge layer

### Acceptance Criteria
- package surface is stable
- no CLI leakage into binding-facing package
- architecture is ready for later bridge work

---

## Phase 8 — Hardening

### Objective
Make the core production-ready.

### Deliverables
- cleanup of unstable APIs
- improved test coverage
- lifecycle edge-case handling
- better shutdown guarantees
- clearer logging
- doc refresh

### Acceptance Criteria
- repeated start and stop scenarios are safe
- errors are predictable
- docs match code
- no dead scaffolding remains
