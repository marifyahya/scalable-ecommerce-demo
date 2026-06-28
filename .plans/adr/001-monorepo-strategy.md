# ADR-001: Monorepo vs. Polyrepo Strategy

**Status:** Accepted  
**Date:** 2026-06-28  
**Deciders:** Project Architect

---

## Context

This project consists of 8 microservices written in 4 different programming languages. We needed to decide whether to host each service in its own repository (polyrepo) or manage all services in a single unified repository (monorepo).

## Decision

We chose the **Monorepo** approach.

## Rationale

| Factor | Monorepo | Polyrepo |
|---|---|---|
| **Onboarding** | One `git clone`, full system visible | Must clone 8 repos separately |
| **Cross-cutting changes** | Single PR for shared config changes | Requires coordinated PRs across multiple repos |
| **`docker-compose.yml`** | Lives at root, orchestrates everything naturally | Would need a separate "orchestration" repo |
| **Learning context** | Much easier to reason about the system as a whole | Creates cognitive overhead irrelevant to learning goals |

## Consequences

**Positive:**
- A single `git clone` and `docker-compose up` gets a new developer running the full system in minutes.
- Shared infrastructure configs (Prometheus, Loki) have one obvious home.

**Negative / Trade-offs:**
- In a real production system with dozens of teams, a monorepo requires investment in tooling (Nx, Turborepo, Bazel) to manage CI/CD per-service. We accept this limitation for a learning project.
- Repository size will grow as all services evolve together.
