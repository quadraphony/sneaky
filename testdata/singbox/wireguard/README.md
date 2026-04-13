# WireGuard via sing-box

This directory contains the real local WireGuard fixture used for sing-box coverage verification.

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

## Verification status

Current label: `partially verified`

What is real and completed:
- key material is real
- config fixture is real
- config detection works through the sing-box adapter
- `sneakycli validate` passes
- `sneakycli probe` reaches runtime startup and fails at traffic stage instead of config stage

What is still blocked:
- this local `sing-box 1.8.10` build has no WireGuard inbound
- passwordless sudo is unavailable, so a privileged kernel/userspace peer could not be provisioned on this machine
- a full local peer loopback path is therefore incomplete in the current environment

## Promotion rule

Do not mark WireGuard as `verified` until:
1. a real peer is provisioned
2. the runtime probe succeeds
3. logs are captured for the successful path
