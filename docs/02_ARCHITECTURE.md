# Architecture

## High-Level Architecture

Sneaky Core uses a layered design:

1. public core wrapper
2. config detection and validation
3. protocol adapter selection
4. runtime lifecycle management
5. logging and stats services

## Design Principles

- one unified external API
- internal protocol adapters
- strong separation of concerns
- testability first
- replacement-friendly upstream integration
- no protocol-specific leakage into the public API unless required

## Main Components

### 1. Core Wrapper

This is the public entry point for the rest of the application.

Responsibilities:
- accept config input
- validate startup conditions
- resolve protocol type
- select adapter
- start and stop runtime
- expose stats and logs

### 2. Config Module

Responsibilities:
- parse raw input
- detect format
- validate required fields
- normalize config data where appropriate
- reject ambiguous or invalid input

### 3. Adapter Module

Responsibilities:
- provide protocol-specific start and stop behavior
- isolate engine-specific logic
- expose consistent runtime status to the wrapper

Each adapter must implement the same internal contract.

### 4. Runtime Module

Responsibilities:
- manage active session lifecycle
- guard against double start
- guard against invalid stop calls
- manage cancellation and cleanup
- ensure clean shutdown behavior

### 5. Observability Module

Responsibilities:
- structured logs
- connection status
- traffic stats
- error reporting surfaces for CLI and later bindings

## Public Core Contract

The public core API must remain small and stable.

Target responsibilities:
- start session
- stop session
- inspect state
- fetch stats
- fetch logs

Do not expand the public API unless there is a strong need.

## Adapter Contract Requirements

Each adapter must define:
- identity
- capability declaration
- config acceptance rules
- startup logic
- shutdown logic
- runtime status reporting

For the SSH adapter family specifically, capability declaration must make it explicit that the adapter:
- provides local SOCKS forwarding
- honors host key checking policy
- honors an explicit known-hosts file when configured

## Config Detection Strategy

Detection must be explicit and deterministic.

Rules:
- do not rely on vague heuristics unless documented
- if detection is ambiguous, return a validation error
- do not silently route one format into another adapter

## Initial Runtime Decision

The first implementation pass must stabilize around the sing-box integration path because it provides the strongest foundation for the initial architecture.

That does not mean the project is sing-box-only forever.
It means the base system must be proven with one real adapter first.

## Mobile Integration Boundary

Mobile bindings are not phase-one deliverables.

Architecture must remain binding-friendly, but implementation priority is:
1. stable core
2. CLI verification
3. binding preparation
4. frontend integration later
