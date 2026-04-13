# Sneaky Core — AI Dev Rules

These rules are mandatory.

## 1. Delivery Rules

- Work strictly in phases.
- Do not jump ahead.
- Do not mix unfinished work from future phases into the current phase.
- Each phase must end in a buildable and testable state.
- Each completed phase must be committed before moving to the next.

## 2. No Guesswork

- Do not invent APIs.
- Do not invent protocol capabilities.
- Do not claim support for a protocol unless the implementation path is explicitly defined in this repo.
- Do not create placeholder handlers that pretend to work.
- Do not create TODO-heavy stubs as if they are finished features.

## 3. No Placeholders

Forbidden:
- fake repos
- fake URLs
- fake package names
- fake bindings
- fake mobile integration code
- placeholder config schemas
- sample values presented as production values

If something is not yet finalized, document it clearly instead of faking it.

## 4. Architecture Rules

- Keep the core modular.
- Use adapter-based protocol design.
- Keep upstream integration isolated.
- Do not tightly couple protocol adapters to CLI code.
- Do not tightly couple config detection to tunnel execution.
- Do not put business logic into `main.go`.

## 5. Source of Truth

The following files are authoritative:
- `docs/01_PRODUCT_SCOPE.md`
- `docs/02_ARCHITECTURE.md`
- `docs/03_PROJECT_STRUCTURE.md`
- `docs/04_PHASE_PLAN.md`
- `docs/05_PROTOCOL_MATRIX.md`

If implementation conflicts with these docs, update the docs first, then implement.

## 6. Upstream Modification Policy

- Do not directly modify upstream engine source unless explicitly approved.
- Prefer wrapper, adapter, extension, or integration layers.
- Keep upstream dependencies replaceable.

## 7. Testing Rules

Every phase must include:
- compile verification
- unit tests where applicable
- manual CLI verification notes
- regression checks for already completed behavior

Do not proceed if:
- build is broken
- tests fail
- logs are unclear
- interfaces are unstable

## 8. Logging Rules

- Use structured logs
- No noisy debug spam in production paths
- No hidden failures
- Every tunnel start failure must produce a useful error
- Every config validation failure must be explicit

## 9. Error Handling Rules

- Fail early on invalid config
- Return typed errors where practical
- Never swallow errors
- Never silently fallback to another protocol without explicit design

## 10. Code Quality Rules

- Keep files focused
- Keep packages cohesive
- Avoid giant utility files
- Avoid unnecessary abstractions
- Avoid premature optimization
- Prefer readable code over clever code

## 11. Scope Control Rules

Do not add:
- frontend code
- payment code
- auth systems
- dashboards
- telemetry services
- cloud sync
- update checkers

unless those are explicitly added to scope later.

## 12. Phase Completion Standard

A phase is complete only when:
- its defined deliverables are implemented
- code builds cleanly
- tests pass
- docs are updated
- commit-ready changes are isolated and understandable
