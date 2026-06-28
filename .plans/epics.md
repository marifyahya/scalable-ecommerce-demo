# Product Backlog (Epics & Tasks)

This document serves as our Jira / Trello / Kanban board replacement. Here, the entire *Development Roadmap* from `blueprint.md` is broken down into sequential **Epics** and **Tasks**.

Completed tasks will be marked with `[x]`.

---

## Epic 1: Foundation & Infrastructure
**Goal:** Build the monorepo foundation and spin up all local databases/brokers.
- [x] **TASK-101:** Create the monorepo folder structure for 7 services.
- [x] **TASK-102:** Configure `docker-compose.yml` for PostgreSQL (4 separate databases).
- [x] **TASK-103:** Add Cassandra, Redis, RabbitMQ, and OpenSearch to Docker.
- [x] **TASK-104:** Add Observability stack (Prometheus, Loki, Grafana) to Docker.

---

## Epic 2: Identity & Edge
**Goal:** Secure the entry point and manage JWT authentication.
- [x] **TASK-201:** Initialize API Gateway (Go/Fiber).
- [x] **TASK-202:** Implement Redis Rate Limiter in Gateway (with bypass support).
- [ ] **TASK-203:** Initialize User Service (Node.js/NestJS), `.env` config, & DB connection.
- [ ] **TASK-204:** Implement Registration endpoint (`POST /register`) & bcrypt.
- [ ] **TASK-205:** Implement Login endpoint (`POST /login`) & JWT Generation.
- [ ] **TASK-206:** Implement Stateless JWT Validation & Header Injection in Gateway.

---

## Epic 3: Catalog & Search (CQRS)
**Goal:** Manage book data and synchronize automatically with the search engine.
- [ ] **TASK-301:** Initialize Product Service (Go/Gin), `.env` config, & DB connection.
- [ ] **TASK-302:** Implement Book Catalog CRUD (Admin).
- [ ] **TASK-303:** Implement Public Catalog endpoints (`GET /products` & `GET /products/:id`).
- [ ] **TASK-304:** Initialize Product Search Service (Python/FastAPI) & OpenSearch connection.
- [ ] **TASK-305:** Implement RabbitMQ Publisher in Product Service (`ProductUpdated`).
- [ ] **TASK-306:** Implement RabbitMQ Consumer in Search Service (Indexing to OpenSearch).
- [ ] **TASK-307:** Implement Search endpoint (`GET /search`).

---

## Epic 4: Shopping Cart
**Goal:** Create a resilient, persistent shopping cart using NoSQL.
- [ ] **TASK-401:** Initialize Cart Service (Node.js/Express) & Cassandra connection.
- [ ] **TASK-402:** Implement Upsert logic (Add to Cart / Update Qty).
- [ ] **TASK-403:** Implement Get Cart and Clear Cart endpoints.

---

## Epic 5: Transactions & Soft-Booking
**Goal:** The heart of e-commerce: stock validation, locking, and order creation.
- [ ] **TASK-501:** Initialize Order Service (PHP/Laravel), `.env` config, & DB connection.
- [ ] **TASK-502:** Implement HTTP Client for stock validation to Product Service (Sync).
- [ ] **TASK-503:** Implement Redis Soft-Booking logic using `SETNX` (15-minute TTL).
- [ ] **TASK-504:** Implement Checkout endpoint (`POST /checkout`) and trigger Cart clearing.
- [ ] **TASK-505:** Implement Cron Job (Laravel Scheduler) to sweep `Expired` orders.

---

## Epic 6: Payments & Async Events
**Goal:** Accept payments and trigger background chain reactions.
- [ ] **TASK-601:** Initialize Payment Service (Node.js/Express) & `.env` config.
- [ ] **TASK-602:** Implement Xendit API integration (Create Invoice).
- [ ] **TASK-603:** Implement Xendit Webhook Callback validation.
- [ ] **TASK-604:** Implement RabbitMQ Publisher (`OrderPaid`).
- [ ] **TASK-605:** Implement RabbitMQ Consumer in Order Service (Update status to PAID).
- [ ] **TASK-606:** Initialize Notification Service (Go) without an HTTP server.
- [ ] **TASK-607:** Implement RabbitMQ Consumer & SMTP Email delivery.

---

## Epic 7: Observability & Load Testing
**Goal:** Monitor system health and push it to its limits.
- [ ] **TASK-701:** Standardize JSON Logging format across all 7 services.
- [ ] **TASK-702:** Setup Grafana Dashboard (Reading Prometheus metrics & Loki logs).
- [ ] **TASK-703:** Write k6 script for Load Testing (Normal Traffic).
- [ ] **TASK-704:** Write k6 script for Spike Testing (Flash Sale using Bypass Header).

---

## Epic 8: Frontend Web (Svelte SPA)
**Goal:** Build a blazing fast user interface without a Virtual DOM.
- [ ] **TASK-801:** Initialize Svelte + Vite + Tailwind CSS project (`frontend-web/`).
- [ ] **TASK-802:** Configure HTTP Client (Axios/Fetch) to hit API Gateway.
- [ ] **TASK-803:** Implement Catalog Page (Book Grid).
- [ ] **TASK-804:** Implement Shopping Cart Page.
- [ ] **TASK-805:** Implement Checkout Page (Integrated with simulated Xendit).
- [ ] **TASK-806:** Implement Order History Page (Calling API Stitching BFF).
