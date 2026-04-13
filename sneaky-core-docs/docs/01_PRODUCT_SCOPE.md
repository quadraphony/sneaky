# Product Scope

## Product Name

Sneaky Core

## Product Type

A modular VPN and tunneling core built in Go.

## Purpose

Sneaky Core exists to provide a single, stable interface for managing multiple connection protocols and tunnel engines without forcing the consuming app to know the internal differences between them.

## Primary Goal

Build a reliable and extensible core that can later be embedded into mobile and other client applications.

## Non-Goals

This repository is not building:
- a Flutter app
- a customer UI
- a billing system
- a website
- account login
- user profile features
- server marketplace
- subscription sales
- analytics dashboards

## Initial Problem Statement

Most tunneling ecosystems are fragmented:
- some configs belong to sing-box-style engines
- some rely on SSH-based transport variations
- some use OpenVPN-style formats
- some require different runtime handling and validation logic

Sneaky Core must hide that complexity behind one core API.

## Core Responsibilities

Sneaky Core must handle:
- config input
- config detection
- config validation
- adapter selection
- tunnel start and stop
- lifecycle management
- logs
- stats
- testable CLI execution

## Success Criteria

The project succeeds when:
1. a single public core interface can start and stop supported adapters
2. config detection is deterministic and testable
3. logs and stats are exposed consistently
4. the codebase is modular enough to add future adapters without refactoring the whole system
5. a CLI can be used to validate the lifecycle outside of any frontend

## Initial Delivery Strategy

Do not start by chasing every protocol.

Start by building:
- the core contract
- the adapter contract
- config detection
- one working adapter foundation
- lifecycle and observability

Then expand protocol coverage phase by phase.
