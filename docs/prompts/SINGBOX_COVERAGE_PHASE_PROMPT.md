# Sneaky Core — Strict AI Dev Prompt for Sing-box Coverage Phase

You are working on **Sneaky Core**.

This phase is **not** about adding new adapter families.
This phase is about **fully verifying and documenting protocol coverage through the existing sing-box adapter**.

You must follow repo rules strictly.

---

## NON-NEGOTIABLE RULES

1. Do not add new standalone adapters in this phase.
2. Do not claim support for any protocol unless it is actually verified.
3. Do not say "all sing-box protocols supported".
4. Do not rely on assumptions from docs alone.
5. Do not fake test results.
6. Do not add placeholders pretending to be production-ready implementations.
7. If something cannot be verified in the current environment, mark it clearly as:
   - blocked
   - not implemented
   - partially verified

---

## PHASE OBJECTIVE

Expand and verify **real protocol coverage under the sing-box adapter**.

This means:
- adding config fixtures
- validating config detection
- validating adapter routing
- running probe or runtime verification where feasible
- updating docs to reflect only what is true

---

## SCOPE FOR THIS PHASE

### Mandatory first batch
- VLESS
- VMess
- Trojan
- Shadowsocks

### Second batch only after the first batch is stable
- WireGuard via sing-box
- Hysteria2
- TUIC
- SSH via sing-box, only if current runtime and binary support can be verified cleanly

### Explicitly out of scope for this phase
- standalone OpenVPN adapter
- standalone WireGuard adapter
- TrustTunnel
- Hiddify-specific adapter
- mobile bindings
- frontend work

---

## REQUIRED OUTPUTS

### 1. Protocol matrix update
Create or update a repo file that clearly marks each protocol as one of:
- verified
- partially verified
- blocked
- not implemented

### 2. Test fixtures
Add config fixtures for each protocol under test.

### 3. Validation coverage
Ensure config detection and adapter selection can recognize and route supported sing-box configs correctly.

### 4. Real verification
For each protocol, provide:
- exact config files used
- exact command run
- observed result
- whether validation passed
- whether runtime/probe passed
- relevant logs/errors

### 5. Final truth table
At the end of the phase, return a concise matrix with:
- protocol
- transport
- detection status
- config validity
- runtime status
- final label

---

## REQUIRED EXECUTION ORDER

### Step 1
Update docs so the repo clearly distinguishes:
- sing-box adapter family
- protocol coverage within sing-box
- future standalone adapters not yet implemented

### Step 2
Add or update fixtures for:
- VLESS
- VMess
- Trojan
- Shadowsocks

### Step 3
Run detection, validation, and runtime/probe verification for that first batch.

### Step 4
Only after first-batch stability, attempt:
- WireGuard via sing-box
- Hysteria2
- TUIC

### Step 5
Update the protocol matrix with actual results only.

---

## VALIDATION STANDARD

A protocol may be marked **verified** only if all of the following are true:
1. a valid sing-box config fixture exists
2. the current adapter detects/routes it correctly
3. config validation passes
4. runtime or probe verification succeeds
5. logs are clean enough to diagnose failures

If any one of those is missing, do not label it verified.

---

## WORDING RULE

Forbidden wording:
- supported
- should work
- probably works
- future-ready
- basically done

Required wording:
- verified
- partially verified
- blocked
- not implemented

---

## OUTPUT FORMAT

Return:
1. file tree of changed files
2. full contents of each changed file
3. commands run
4. per-protocol evidence
5. final updated protocol matrix

Proceed now with documentation correction first, then verified sing-box protocol coverage.
