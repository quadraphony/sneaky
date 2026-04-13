# Sneaky Core — Protocol Matrix

This file is the source of truth for protocol claims in this repository.

## Status Labels

- **verified** = config exists, detection works, validation passes, runtime/probe passes
- **partially verified** = some validation exists, but full runtime verification is incomplete
- **blocked** = cannot be completed in the current environment
- **not implemented** = no real implementation work exists yet

## Adapter Families

### Implemented adapter families
- sing-box adapter
- ssh direct SOCKS adapter

### Future adapter families
- standalone OpenVPN adapter
- standalone WireGuard adapter
- TrustTunnel adapter
- Hiddify-specific adapter

These future adapter families are **not implemented** unless explicitly stated otherwise.

---

## sing-box Coverage Matrix

| Protocol / Capability | Category | Repo target status | Verification label | Notes |
|---|---|---|---|---|
| VLESS | proxy protocol | verified in coverage phase | verified | local loopback runtime and probe passed |
| VMess | proxy protocol | verified in coverage phase | verified | local loopback runtime and probe passed |
| Trojan | proxy protocol | verified in coverage phase | verified | local TLS loopback runtime and probe passed |
| Shadowsocks | proxy protocol | verified in coverage phase | verified | local loopback runtime and probe passed |
| WireGuard via sing-box | proxy / tunnel | test after first batch | blocked | key generation and peer wiring not yet automated |
| Hysteria2 | proxy protocol | test after first batch | blocked | QUIC/TLS verification path not completed in this repo |
| TUIC | proxy protocol | test after first batch | blocked | QUIC/TLS verification path not completed in this repo |
| SSH via sing-box | proxy protocol | optional after first batch | blocked | repo has standalone SSH adapter; sing-box SSH path not verified |
| Naive | proxy protocol | deferred | not implemented | phase-later candidate |
| ShadowTLS | proxy protocol | deferred | not implemented | phase-later candidate |
| AnyTLS | proxy protocol | deferred | not implemented | phase-later candidate |
| Hysteria | proxy protocol | deferred | not implemented | prefer Hysteria2 first |
| Tor | outbound capability | deferred | not implemented | environment dependent |
| SOCKS | utility outbound/inbound | utility path only | partially verified | used in runtime probe path, not primary protocol target |
| HTTP CONNECT | outbound type | deferred | not implemented | can be added later as utility path |
| DNS outbound | internal capability | out of verification focus | not implemented | not part of primary coverage phase |
| Selector | utility outbound | out of verification focus | not implemented | management utility, not primary protocol |
| URLTest | utility outbound | out of verification focus | not implemented | management utility, not primary protocol |
| Direct | utility outbound | baseline utility | partially verified | used by loopback servers |
| Block | utility outbound | baseline utility | partially verified | not a proxy protocol target |

---

## Verification Promotion Rule

No row may move to **verified** unless:
1. fixture exists
2. detection succeeds
3. validation succeeds
4. runtime/probe succeeds
5. logs are captured

---

## Current Strategy

This repo verified, in order:
1. VLESS
2. VMess
3. Trojan
4. Shadowsocks

The next sing-box candidates remain:
5. WireGuard via sing-box
6. Hysteria2
7. TUIC
