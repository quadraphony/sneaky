# WireGuard via sing-box — deferred fixture

This fixture is intentionally deferred until key generation is handled inside the repo test workflow.

## Why deferred

WireGuard in sing-box requires:
- local interface addresses
- private key
- peer public key
- peer endpoint details

These are required fields and should not be faked.

## Required prep

Generate keys with either:
- `wg genkey` and `wg pubkey`
- or `sing-box generate wg-keypair`

## Promotion rule

Do not mark WireGuard as verified until:
1. keys are generated correctly
2. local_address is assigned correctly
3. peer definitions are valid
4. runtime/probe passes
5. logs are captured
