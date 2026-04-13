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
