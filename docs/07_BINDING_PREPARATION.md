# Binding Preparation

## Objective

Prepare the Go core for later Android and iOS bridge layers without exposing internal implementation packages.

## Stable Public Surface

The binding-facing package is `pkg/sneaky`.

It currently exposes:
- manager construction
- start and stop lifecycle control
- state snapshot access
- structured log access
- manager-level stats access
- config inspection helpers for file and raw bytes
- typed public adapter identifiers
- explicit public lifecycle state helpers

Bridge layers must not import `internal/*` packages directly.

## Binding Safety Notes

- lifecycle control is synchronous and explicit
- public types are plain Go values without CLI dependencies
- config inspection is available without constructing CLI commands
- adapter registration remains internal so mobile bridges do not manage engine wiring
- adapter selection does not require raw magic strings in bridge code

## CLI Boundary

The CLI remains a development tool only.

No CLI-only types are part of `pkg/sneaky`.

## Future Bridge Work

Later bridge layers should wrap:
- `Manager.Start`
- `Manager.Stop`
- `Manager.Snapshot`
- `Manager.Logs`
- `Manager.Stats`
- `InspectConfigPath`
- `InspectConfigBytes`
