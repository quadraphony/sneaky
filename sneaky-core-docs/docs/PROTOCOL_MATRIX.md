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
| WireGuard via sing-box | proxy / tunnel | migration follow-up only | blocked | legacy outbound fixture is pinned to pre-1.11 behavior; current repo runtime uses sing-box 1.13.7 where WireGuard moved to endpoint configuration |
| Hysteria2 | proxy protocol | verified in coverage phase | verified | local QUIC/TLS loopback runtime and probe passed |
| TUIC | proxy protocol | verified in coverage phase | verified | local QUIC/TLS loopback runtime and probe passed |
| SSH via sing-box | proxy protocol | verified in coverage phase | verified | local user-space `sshd` runtime and probe passed |
| Naive | proxy protocol | verified in coverage phase | verified | local loopback runtime and probe passed with bundled sing-box 1.13.7 |
| ShadowTLS | proxy protocol | verified in coverage phase | verified | local loopback runtime and probe passed with chained ShadowTLS and Shadowsocks outbounds |
| AnyTLS | proxy protocol | verified in coverage phase | verified | local loopback runtime and probe passed with bundled sing-box 1.13.7 |
| Hysteria | proxy protocol | verified in coverage phase | verified | local QUIC/TLS loopback runtime and probe passed |
| Tor | outbound capability | verified in coverage phase | verified | local bundled Tor executable bootstrapped and probe passed |
| SOCKS | utility outbound/inbound | utility path only | partially verified | used in runtime probe path, not primary protocol target |
| HTTP CONNECT | utility outbound/inbound | verified utility path | verified | local loopback runtime and probe passed |
| DNS outbound | internal capability | out of verification focus | partially verified | utility fixture validates with `sing-box check` |
| Selector | utility outbound | out of verification focus | not implemented | management utility, not primary protocol |
| Selector | utility outbound | out of verification focus | partially verified | utility fixture validates with `sing-box check` |
| URLTest | utility outbound | out of verification focus | partially verified | utility fixture validates with `sing-box check` |
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
5. Hysteria2
6. TUIC
7. Hysteria
8. HTTP CONNECT
9. ShadowTLS
10. AnyTLS
11. Naive
12. SSH via sing-box
13. Tor

The next sing-box candidates remain:
14. WireGuard via sing-box migration
