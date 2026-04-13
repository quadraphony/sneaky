# Protocol Matrix

This file defines what is actually in scope and how protocol support must be treated.

## Important Rule

A protocol is not considered supported just because it is listed here.

Support status must be one of:
- planned
- under evaluation
- active
- blocked
- out of scope

## Current Matrix

| Protocol Family | Status | Implementation Path Defined | Notes |
|---|---|---:|---|
| sing-box based configs | active foundation target | yes | first real adapter path |
| OpenVPN | planned | not yet finalized | only after adapter architecture stabilizes |
| SSH-based tunneling | planned | not yet finalized | must be specified precisely before implementation |
| WireGuard | planned | not yet finalized | do not promise until adapter path is selected |
| TrustTunnel | under evaluation | no | do not implement until the source and integration path are fully locked |
| Hiddify-based path | under evaluation | no | only adopt if the integration approach is explicitly approved |

## Protocol Expansion Rule

Before any protocol moves from `planned` or `under evaluation` to `active`, the following must exist:
1. implementation source identified
2. adapter contract mapped
3. validation strategy defined
4. runtime lifecycle defined
5. test strategy defined

## SSH Rule

“SSH-based support” is too broad unless broken down precisely.

Before implementation, define the exact supported forms, for example:
- direct SSH tunnel
- SSH over TLS
- SSH over WebSocket
- payload-based HTTP header injection
- DNS-based transport

Do not implement “SSH support” as a vague bucket.

## OpenVPN Rule

OpenVPN support requires:
- config identification rules
- runtime invocation strategy
- lifecycle ownership model
- log extraction strategy
- shutdown strategy

## TrustTunnel Rule

TrustTunnel must remain under evaluation unless:
- the codebase source is locked
- ownership and maintenance risk are known
- build integration is defined
- runtime contract is defined

No fake TrustTunnel support is allowed.
