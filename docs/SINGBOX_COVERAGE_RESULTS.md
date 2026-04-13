# Sing-box Coverage Results

This document records the exact loopback verification evidence used to promote sing-box protocol rows to **verified**.

## Verification Method

For each protocol:
1. validate server and client fixtures with `sing-box check`
2. start the loopback server fixture with `sing-box run`
3. run `sneakycli probe` against the client fixture
4. record the observed probe result

## Verified Batch

### VLESS
- Server fixture: `testdata/singbox/vless/server.json`
- Client fixture: `testdata/singbox/vless/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/vless/server.json`
  - `sing-box check -c testdata/singbox/vless/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/vless/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/vless/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### VMess
- Server fixture: `testdata/singbox/vmess/server.json`
- Client fixture: `testdata/singbox/vmess/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/vmess/server.json`
  - `sing-box check -c testdata/singbox/vmess/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/vmess/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/vmess/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### Trojan
- Server fixture: `testdata/singbox/trojan/server.json`
- Client fixture: `testdata/singbox/trojan/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/trojan/server.json`
  - `sing-box check -c testdata/singbox/trojan/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/trojan/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/trojan/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### Shadowsocks
- Server fixture: `testdata/singbox/shadowsocks/server.json`
- Client fixture: `testdata/singbox/shadowsocks/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/shadowsocks/server.json`
  - `sing-box check -c testdata/singbox/shadowsocks/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/shadowsocks/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/shadowsocks/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### Hysteria2
- Server fixture: `testdata/singbox/hysteria2/server.json`
- Client fixture: `testdata/singbox/hysteria2/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/hysteria2/server.json`
  - `sing-box check -c testdata/singbox/hysteria2/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/hysteria2/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/hysteria2/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### TUIC
- Server fixture: `testdata/singbox/tuic/server.json`
- Client fixture: `testdata/singbox/tuic/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/tuic/server.json`
  - `sing-box check -c testdata/singbox/tuic/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/tuic/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/tuic/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### Hysteria
- Server fixture: `testdata/singbox/hysteria/server.json`
- Client fixture: `testdata/singbox/hysteria/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/hysteria/server.json`
  - `sing-box check -c testdata/singbox/hysteria/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/hysteria/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/hysteria/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### HTTP CONNECT
- Server fixture: `testdata/singbox/http/server.json`
- Client fixture: `testdata/singbox/http/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/http/server.json`
  - `sing-box check -c testdata/singbox/http/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/http/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/http/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### ShadowTLS
- Server fixture: `testdata/singbox/shadowtls/server.json`
- Client fixture: `testdata/singbox/shadowtls/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/shadowtls/server.json`
  - `sing-box check -c testdata/singbox/shadowtls/client.json`
  - `go run ./cmd/sneakycli validate testdata/singbox/shadowtls/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/shadowtls/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/shadowtls/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - server logs showed `inbound/shadowsocks[ss-in]: inbound connection to example.com:443`
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### AnyTLS
- Server fixture: `testdata/singbox/anytls/server.json`
- Client fixture: `testdata/singbox/anytls/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/anytls/server.json`
  - `sing-box check -c testdata/singbox/anytls/client.json`
  - `go run ./cmd/sneakycli validate testdata/singbox/anytls/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/anytls/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/anytls/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - server logs showed `inbound/anytls[anytls-in]: [local-test] inbound connection to example.com:443`
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### Naive
- Server fixture: `testdata/singbox/naive/server.json`
- Client fixture: `testdata/singbox/naive/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/naive/server.json`
  - `sing-box check -c testdata/singbox/naive/client.json`
  - `go run ./cmd/sneakycli validate testdata/singbox/naive/client.json`
- Runtime command:
  - `sing-box run -c testdata/singbox/naive/server.json`
  - `go run ./cmd/sneakycli probe testdata/singbox/naive/client.json https://example.com`
- Observed result:
  - loopback server started successfully
  - server logs showed `inbound/naive[naive-in]: [local-test] inbound connection to example.com:443`
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### SSH via sing-box
- Client fixture: `testdata/singbox/ssh/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/ssh/client.json`
  - `go run ./cmd/sneakycli validate testdata/singbox/ssh/client.json`
- Runtime command:
  - start local user-space `sshd` with repo test keys on `127.0.0.1:22322`
  - `go run ./cmd/sneakycli probe testdata/singbox/ssh/client.json https://example.com`
- Observed result:
  - local `sshd` accepted the repo test key
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

### Tor
- Client fixture: `testdata/singbox/tor/client.json`
- Validation command:
  - `sing-box check -c testdata/singbox/tor/client.json`
  - `go run ./cmd/sneakycli validate testdata/singbox/tor/client.json`
- Runtime command:
  - `go run ./cmd/sneakycli probe testdata/singbox/tor/client.json https://example.com`
- Observed result:
  - local bundled Tor executable bootstrapped to 100%
  - `sneakycli probe` returned `probe ok adapter=singbox ... status=200 ...`
- Final label: `verified`

## Blocked Rows

### WireGuard
- Fixture:
  - `testdata/singbox/wireguard/client.json`
- Historical validation context:
  - originally checked under local `sing-box 1.8.10`
- Observed result:
  - the legacy fixture no longer matches the repo's bundled `sing-box 1.13.7` runtime because WireGuard moved to endpoint configuration in sing-box 1.11.0+
  - the old partial-verification evidence is retained in `testdata/singbox/wireguard/README.md`
- Final label: `blocked`
