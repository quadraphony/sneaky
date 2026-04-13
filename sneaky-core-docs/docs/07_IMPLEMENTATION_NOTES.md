# Implementation Notes

## Initial Technical Decision

The first real adapter must be sing-box based.

Reason:
- it gives a realistic engine-backed foundation
- it allows the core contracts to be tested against a real runtime path
- it reduces early architectural guesswork

## Why Not Start With Everything

Starting with multiple adapters at once creates:
- unstable contracts
- unclear validation rules
- harder debugging
- weak test quality
- fake abstraction pressure

## What “Modular” Means Here

Modular does not mean:
- dumping every protocol into one manager
- switching by giant `if else` chains forever
- mixing runtime code and parsing code

Modular means:
- small stable contracts
- isolated adapters
- controlled registration
- explicit capability boundaries

## Documentation Rule

Whenever implementation changes architecture, update docs first or in the same change set.
