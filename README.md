# Tukubuku

**Status: Under Development (Work In Progress)**

A polyglot e-commerce microservices system built for learning and portfolio purposes. The architecture applies enterprise-grade patterns — event-driven communication, CQRS, distributed data stores, and a full observability stack — to a small bookstore domain.

## Architecture Overview

| Property | Value |
|---|---|
| Strategy | Monorepo |
| Architecture | Microservices (Polyglot) |
| Orchestration | Docker Compose |
| Services | 8 (1 API Gateway + 7 domain services) |
| Async Messaging | RabbitMQ |
| Observability | Prometheus + Loki + Grafana |

For detailed architecture decisions, service boundaries, and communication rules, see [.plans/blueprint.md](./.plans/blueprint.md).

## Repository Structure

```
tukubuku/
├── frontend-web/            # TypeScript (Svelte + Vite) — SPA user interface
├── api-gateway/             # Go — entry point, routing, JWT validation, rate limiting
├── user-service/            # Node.js (NestJS) — authentication & profile management
├── product-service/         # Go — master book catalog (source of truth)
├── product-search-service/  # Python (FastAPI) — OpenSearch indexing & full-text search
├── cart-service/            # Node.js (Express) — persistent cart (Cassandra)
├── order-service/           # PHP (Laravel) — checkout & order lifecycle
├── payment-service/         # Node.js (Express) — Xendit webhook integration
├── notification-service/    # Go/Python — async email delivery worker
├── load-tests/              # k6 scripts for load, spike, and stress testing
├── docker/                  # Infrastructure configs (Prometheus, Loki, Grafana, DB init)
└── docker-compose.yml       # Single orchestration file for the full ecosystem
```

## Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (macOS/Windows) or Docker Engine (Linux)
- Git

### Local Port Configuration

Port mappings are not defined in `docker-compose.yml` by default. To expose service ports to your host machine, copy the provided example override file and adjust it as needed:

```bash
cp docker-compose.override.yml.example docker-compose.override.yml
```

Docker Compose automatically merges `docker-compose.override.yml` on startup. The file is gitignored, so each developer can customize it independently.

### Running the Infrastructure

From the project root, start all services in detached mode:

```bash
docker-compose up -d
```

This will automatically provision the following containers:

| Container | Description |
|---|---|
| `postgres_db` | PostgreSQL 15 with 4 isolated databases pre-created |
| `cassandra_db` | Apache Cassandra 4.1 for the cart service |
| `redis_cache` | Redis 7 for rate limiting and soft-booking locks |
| `rabbitmq_broker` | RabbitMQ 3 with management plugin |
| `opensearch_node` | OpenSearch 2.11 for full-text product search |
| `prometheus` | Metrics collection |
| `loki` | Log aggregation |
| `grafana` | Unified observability dashboard |

### Local Access

After startup, the following interfaces are available on your host machine:

| Service | URL | Credentials |
|---|---|---|
| RabbitMQ Management | http://localhost:15672 | `guest` / `guest` |
| Grafana | http://localhost:3000 | `admin` / `admin` |
| OpenSearch | http://localhost:9200 | — |
| Prometheus | http://localhost:9090 | — |

## Documentation

All technical documentation is located in the [.plans/](./.plans/) directory.

| Document | Description |
|---|---|
| [prd.md](./.plans/prd.md) | Product requirements, user journeys, and business rules |
| [blueprint.md](./.plans/blueprint.md) | Architecture design, tech stack decisions, and communication rules |
| [epics.md](./.plans/epics.md) | Development backlog broken down into epics and tasks |
| [adr/](./.plans/adr/) | Architecture Decision Records — context, decision, and trade-offs |
| [services/](./.plans/services/) | Per-service documentation: endpoints, DB schema, and env vars |

## Development Workflow

All active and upcoming work is tracked in [.plans/epics.md](./.plans/epics.md). Each task is organized under an epic and marked upon completion.

To pick up work, refer to that file for the next unchecked task in the current active epic.
