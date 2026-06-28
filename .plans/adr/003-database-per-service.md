# ADR-003: Database-per-Service Pattern

**Status:** Accepted  
**Date:** 2026-06-28  
**Deciders:** Project Architect

---

## Context

In a traditional monolith, all services share one database and use SQL JOINs freely. In microservices, we needed to decide whether to allow services to share a database schema or enforce strict isolation.

## Decision

Enforce the **Database-per-Service** pattern. No service may access another service's database directly.

## Rationale

Sharing a database creates hidden coupling:
- Schema changes in one service can silently break another.
- A single database becomes the system's bottleneck and single point of failure.
- It undermines the primary microservices goal: independent deployability.

With Database-per-Service, if the Cart Service's Cassandra cluster goes down, the Order Service's PostgreSQL remains unaffected and checkouts already in progress continue normally.

## How Cross-Service Data Retrieval Works

Since SQL JOINs across service boundaries are forbidden, two alternatives are used:

1. **Synchronous HTTP (for critical reads):** Order Service calls Product Service via HTTP GET at checkout time to fetch live price and stock data.
2. **API Stitching at Gateway (for UI composition):** The API Gateway fetches from multiple services and assembles the final response (e.g., Order history enriched with book titles from Product Service).

## Consequences

**Positive:**
- True service independence — each service can be scaled, replaced, or failed in isolation.
- Different services can use different database technologies optimized for their access patterns (PostgreSQL for relational, Cassandra for write-heavy, Redis for ephemeral state).

**Negative / Trade-offs:**
- No cross-service JOIN queries. Data that used to be one query is now 2+ network calls.
- Eventual consistency is the norm — OpenSearch may briefly lag behind PostgreSQL after a product update.
- Requires disciplined API contract design between services.
