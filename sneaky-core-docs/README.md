# Sneaky Core

A modular Go-based VPN and tunneling core designed to unify multiple connection protocols behind one stable interface.

This repository is **core-only**. It does **not** include the mobile or desktop frontend.  
Frontend integration will happen later in a separate repository.

## Project Objective

Sneaky Core must provide:

- a single entry point for starting and stopping tunnels
- protocol-specific adapters behind one unified interface
- configuration detection and validation
- structured logging
- traffic statistics
- clean integration points for Android and iOS bindings later
- a CLI for development, testing, and debugging

## Current Product Boundary

This repository is responsible for:

- core architecture
- protocol adapter system
- tunnel lifecycle management
- config parsing and validation
- logging and stats
- CLI testing entrypoint
- mobile binding preparation

This repository is **not** responsible for:

- Flutter UI
- mobile screen design
- subscription logic
- payments
- marketing site
- onboarding UX

## Engineering Principles

- no placeholders
- no fake implementations
- no “coming soon” code paths
- no protocol listed unless its implementation path is defined
- no direct upstream source edits unless explicitly approved
- phased delivery only
- every phase must compile and pass tests before the next begins

## Initial Technical Direction

- Language: Go
- Primary base engine: sing-box
- Architecture style: modular adapter-based core
- Testing entrypoint: CLI
- Mobile readiness: planned through bindings after core stabilizes

## Initial Supported Scope

Phase 1 support focuses on:

- unified core wrapper
- config detection
- sing-box adapter foundation
- structured logs
- connection stats
- CLI lifecycle commands

Additional protocol support will be added only in later phases once each has:
1. a confirmed implementation path
2. a defined adapter contract
3. a test strategy

## Repository Structure

See:
- `docs/01_PRODUCT_SCOPE.md`
- `docs/02_ARCHITECTURE.md`
- `docs/03_PROJECT_STRUCTURE.md`
- `docs/04_PHASE_PLAN.md`
- `docs/05_PROTOCOL_MATRIX.md`
- `docs/AI_DEV_RULES.md`

## Development Rule

Do not start broad protocol expansion before the base architecture is stable.

The first milestone is not “support everything”.
The first milestone is “build a correct, extensible core”.

## Status

Planning and architecture phase.
Implementation must follow the defined phase plan in `docs/04_PHASE_PLAN.md`.
