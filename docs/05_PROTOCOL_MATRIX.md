# Protocol Matrix

This file defines what is actually in scope and how protocol support must be treated.

## Important Rule

A protocol is not considered supported just because it is listed here.

Detailed sing-box protocol coverage lives in `docs/PROTOCOL_MATRIX.md`.

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
| SSH direct SOCKS tunnel | active | yes | implemented as OpenSSH dynamic forwarding via `ssh -N -D` |
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

Only one SSH form is in scope now:
- direct SSH tunnel with local SOCKS forwarding using OpenSSH `ssh -N -D`

The supported config model is JSON with a top-level `ssh_tunnel` object:
- `host` string, required
- `user` string, required
- `port` integer, optional, defaults to `22`
- `local_socks_port` integer, required
- `identity_file` string, optional
- `known_hosts_file` string, optional
- `strict_host_key_checking` string, optional, defaults to `accept-new`
- `connect_timeout_seconds` integer, optional, defaults to `10`
- `server_alive_interval_seconds` integer, optional, defaults to `30`
- `server_alive_count_max` integer, optional, defaults to `3`

The following SSH forms remain out of scope until separately specified:
- SSH over TLS
- SSH over WebSocket
- payload/header injection modes
- DNS-based transports

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
