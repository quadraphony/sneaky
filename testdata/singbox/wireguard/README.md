# WireGuard via sing-box

This directory contains the historical WireGuard fixture used for the older sing-box outbound model.

## Key generation approach

Keys were generated with:

```bash
sing-box generate wg-keypair
sing-box generate wg-keypair
```

Generated keypairs:

Peer A:
- PrivateKey: `cIwzFxjrEqo97GxhZJJTH2gWpbCaTkaheXdbwKVMMH0=`
- PublicKey: `bAJkbEGP8k7LZw8NSFgdQkOBgzQIKsVRjvOXPqinCU4=`

Peer B:
- PrivateKey: `yLGQn69AeVUP5JJm4dJPijk3OlnknRWnasOSkYK1v28=`
- PublicKey: `KCOS4+I08yP0koDSMHyI0siOyZqPp4pUByaKJ8izdDw=`

The checked-in client fixture uses:
- peer A private key
- peer B public key
- local address `10.7.0.2/32`
- endpoint `127.0.0.1:19100`

## Current status

Current label: `blocked`

What remains useful:
- key material is real
- fixture history is real
- the original partial-verification evidence remains documented

What is blocked now:
- the repo now uses bundled `sing-box 1.13.7`
- WireGuard moved to endpoint configuration in sing-box 1.11.0+
- this checked-in file still uses the older outbound-oriented shape and is retained only as migration evidence

## Promotion rule

Do not mark WireGuard as `verified` until:
1. an endpoint-based fixture is added for the current runtime
2. a real peer is provisioned
3. the runtime probe succeeds
4. logs are captured for the successful path
