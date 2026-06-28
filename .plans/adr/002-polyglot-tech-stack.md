# ADR-002: Polyglot Tech Stack

**Status:** Accepted  
**Date:** 2026-06-28  
**Deciders:** Project Architect

---

## Context

This is a learning project. We had the option to standardize on a single language (e.g., all Node.js or all Go) which would minimize context-switching overhead, or to intentionally use multiple languages (polyglot) to simulate a real-world enterprise environment where different teams pick the best tool for their domain.

> **⚠️ Learning Project Note:** A single-language approach would be more practical for a real small-scale bookstore. The polyglot choice is justified *solely* by the educational objective of learning how different runtimes interoperate within the same Docker network.

## Decision

We chose a **Polyglot** approach with the following assignments:

| Service | Language |
|---|---|
| API Gateway | Go |
| User Service | Node.js (NestJS) |
| Product Service | Go |
| Product Search Service | Python (FastAPI) |
| Cart Service | Node.js (Express) |
| Order Service | PHP (Laravel) |
| Payment Service | Node.js (Express) |
| Notification Service | Go or Python |

## Rationale

Each language was chosen for genuine technical reasons within its service's domain:
- **Go** at the network edge (Gateway) and catalog (Product) for raw concurrency and low memory footprint.
- **NestJS** for identity because its opinionated structure is ideal for auth, and the JS ecosystem has the best JWT/bcrypt libraries.
- **Laravel** for order management because its Eloquent ORM, built-in scheduler, and queue integration reduce implementation time for complex business logic.
- **Python** for search sync because of its superior text-processing ecosystem for OpenSearch integration.

## Consequences

**Positive:**
- Directly mirrors how large engineering organizations operate.
- Forces the engineer to understand how services communicate across language boundaries through well-defined contracts (HTTP JSON, RabbitMQ events).
- Demonstrates the value of the Structured Logging JSON standard as the universal language of observability.

**Negative / Trade-offs:**
- Higher cognitive overhead than a single-language stack.
- Each service requires separate Docker image builds and language-specific dependency management.
- A single developer must context-switch between Go, PHP, Python, and Node.js — acceptable for learning, painful for a production team of one.
