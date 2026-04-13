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

## Current Partial Verification

- WireGuard via sing-box

Reason:
- real key material exists
- real outbound fixture exists
- validation passes
- runtime probe reaches startup but fails traffic proof because this machine does not currently provide a real WireGuard peer path for sing-box 1.8.10

## Current Blocked Items

- SSH via sing-box
  - blocked because local `sshd` runtime is not installed
- Naive
  - blocked because local `sing-box 1.8.10` reports `unknown outbound type: naive`
- AnyTLS
  - blocked because local `sing-box 1.8.10` reports `unknown outbound type: anytls`
- Tor
  - blocked because runtime probe failed on this machine

## Current Partial Utility Coverage

- ShadowTLS
  - config type is recognized
  - no runtime/probe fixture implemented yet
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
   - create real local fixture if possible
   - validate config detection and routing
   - run real runtime/probe verification
   - promote only if runtime proof succeeds

2. Reassess Naive and AnyTLS
   - first determine whether the local sing-box runtime must be upgraded
   - do not add fixtures until runtime support is confirmed

3. Reassess SSH via sing-box
   - only if a real `sshd` test path is available

4. Decide whether to close the sing-box coverage phase
   - if remaining rows are environment-blocked or version-blocked, stop expanding claims
   - shift to core hardening and CLI polish

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
