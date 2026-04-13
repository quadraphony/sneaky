# Sneaky Core — Next Steps

This file is the current handoff point for the next development session.

## Current Verified sing-box Coverage

- VLESS
- VMess
- Trojan
- Shadowsocks
- Hysteria2
- TUIC
- Hysteria
- HTTP CONNECT
- ShadowTLS
- AnyTLS
- Naive
- SSH via sing-box
- Tor

## Sing-box Phase Status

- closed
- coverage claims are complete for the current bundled runtime
- do not expand sing-box protocol scope further unless a new product requirement appears

## Deferred Migration Item

- WireGuard via sing-box
  - the historical fixture is kept for evidence, but it targets the pre-1.11 outbound schema
  - the repo now uses bundled `sing-box 1.13.7`, where WireGuard moved to endpoint configuration
  - deferred unless current product scope still needs sing-box-managed WireGuard
  - if revived later, add a new endpoint-based fixture instead of reviving the old outbound path

## Current Partial Utility Coverage

- DNS outbound
  - validation fixture exists
- Selector
  - validation fixture exists
- URLTest
  - validation fixture exists
- SOCKS
  - used by probe path
- Direct
  - used by loopback servers
- Block
  - utility baseline only

## Required Next Order

1. ShadowTLS
   - already verified
   - do not reopen unless a regression appears

2. Phase 1 hardening
   - start the next hardening pass on manager reliability, startup and shutdown behavior, and CLI stability

3. Final public API
   - lock the exported package surface for external callers and future bindings

4. SSH adapter family
   - continue the dedicated SSH adapter path as a first-class family beyond sing-box SSH outbound coverage

5. WireGuard via sing-box
   - deferred migration item
   - only revisit if current product scope still requires it

## Evidence Sources

Use these files first in the next session:
- `docs/PROTOCOL_MATRIX.md`
- `docs/SINGBOX_COVERAGE_RESULTS.md`
- `testdata/singbox/wireguard/README.md`

## Rule Reminder

Do not promote any row to `verified` unless:
1. fixture exists
2. detection succeeds
3. validation succeeds
4. runtime/probe succeeds
5. logs or observed output are captured
