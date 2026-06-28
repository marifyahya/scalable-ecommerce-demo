# 📚 E-Commerce Microservices — Documentation Index

> **⚠️ Learning Project Disclaimer**
> This project is intentionally architected for **learning and portfolio purposes**. The technology stack (polyglot languages, Cassandra, OpenSearch, RabbitMQ, full observability stack) is considered **overkill for its actual business scope** (a simple bookstore). In a real production environment, you would start with a monolith or 2-3 services. The complexity here is deliberate — designed to expose the engineer to enterprise-grade patterns in a controlled setting.

---

## Documentation Map

| Document | Audience | Purpose |
|---|---|---|
| [PRD — Product Requirements](./prd.md) | Everyone | What we're building and why. User journeys, actors, and business rules. |
| [Blueprint — Architecture](./blueprint.md) | Tech Lead / Senior Engineer | How it's built. Tech stack, folder structure, service boundaries, cross-service communication rules, and advanced patterns. |
| [ADR — Architecture Decisions](./adr/) | Senior Engineer / Future Maintainer | Why specific technology choices were made. Context, decision, and trade-offs. |
| [Services — Per-Service Detail](./services/) | Developer implementing a service | Deep-dive per service: endpoints, DB schema, environment variables, and sequence diagrams. |
| [Epics & Tasks — Backlog](./epics.md) | Developer & Project Manager | The Jira/Kanban board breaking down the roadmap into actionable tasks. |

---

## Quick Start

**Onboarding a new team member?** Read in this order:
1. `README.md` (this file) — get oriented
2. `prd.md` — understand the product and user flows
3. `blueprint.md` — understand the overall architecture and communication rules
4. `adr/` — understand *why* key decisions were made
5. `services/<your-service>.md` — understand the specific service you're implementing

---

## Project at a Glance

| Property | Value |
|---|---|
| **Type** | Learning / Portfolio Project |
| **Domain** | E-Commerce (Books only) |
| **Architecture** | Microservices (Polyglot) |
| **Complexity** | ⚠️ Intentionally High (Enterprise-grade patterns) |
| **Repo Strategy** | Monorepo |
| **Orchestration** | Docker Compose (local dev) |
| **Total Services** | 8 (API Gateway + 7 domain services) |

---

## Service Index

| Service | Language | Store | Doc |
|---|---|---|---|
| Frontend Web | TypeScript (Svelte+Vite) | Browser Storage | [frontend-web.md](./services/frontend-web.md) |
| API Gateway | Go (Fiber/Gin) | Redis | [api-gateway.md](./services/api-gateway.md) |
| User Service | Node.js (NestJS) | PostgreSQL | [user-service.md](./services/user-service.md) |
| Product Service | Go (Gin) | PostgreSQL | [product-service.md](./services/product-service.md) |
| Product Search Service | Python (FastAPI) | OpenSearch | [product-search-service.md](./services/product-search-service.md) |
| Cart Service | Node.js (Express) | Cassandra | [cart-service.md](./services/cart-service.md) |
| Order Service | PHP (Laravel) | PostgreSQL + Redis | [order-service.md](./services/order-service.md) |
| Payment Service | Node.js (Express) | PostgreSQL | [payment-service.md](./services/payment-service.md) |
| Notification Service | Go / Python | — | [notification-service.md](./services/notification-service.md) |

---

## ADR Index

| ADR | Decision |
|---|---|
| [ADR-001](./adr/001-monorepo-strategy.md) | Monorepo over Polyrepo |
| [ADR-002](./adr/002-polyglot-tech-stack.md) | Polyglot language strategy |
| [ADR-003](./adr/003-database-per-service.md) | Database-per-Service isolation |
| [ADR-004](./adr/004-rabbitmq-async-communication.md) | RabbitMQ for async communication |
| [ADR-005](./adr/005-stateless-jwt-at-gateway.md) | Stateless JWT validation at Gateway |
| [ADR-006](./adr/006-soft-booking-via-redis.md) | Soft-Booking via Redis for stock reservation |
