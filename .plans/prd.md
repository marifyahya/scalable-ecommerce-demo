# Product Requirements Document (PRD)

**Project:** E-Commerce Microservices  
**Status:** In Design  
**Last Updated:** 2026-06-28

> **⚠️ Learning Project Disclaimer**
> The primary goal of this project is **not** business completeness. It is a deliberate exercise in implementing enterprise-grade microservices architecture patterns: inter-service communication, event-driven design, stateful booking, and system resilience. The chosen tech stack is intentionally more complex than the domain requires.

---

## 1. Objective

Build a small-scale e-commerce platform using microservices architecture standards found in large-scale systems (Shopify, Netflix, Amazon). The focus is on:

- **Inter-service communication** (synchronous vs. asynchronous)
- **State management** (soft-booking, distributed caching)
- **Event-driven architecture** (RabbitMQ)
- **System resilience under load** (rate limiting, network isolation, load testing)

---

## 2. Actors

| Actor | Description |
|---|---|
| **Customer** | End user who browses, adds to cart, and purchases books. |
| **Admin** | System operator who manages the product catalog and order statuses. |

---

## 3. Application Flows

### 3.1 Customer Flow

Customer data is intentionally minimal: `name`, `email`, `password`.

**1. Registration & Login**
- Customer submits credentials via the application. The API Gateway forwards unauthenticated requests directly to **User Service**.
- **User Service** verifies the password against PostgreSQL. On success, it is the **sole authority responsible for generating and signing the JWT token** (containing `user_id` and `role`).
- *Security Note:* The API Gateway and User Service share the same secret/public key, enabling the Gateway to perform **stateless JWT validation** on subsequent requests without calling User Service again.

**2. Product Browsing**
- Customer views the book catalog (served by Product Service via PostgreSQL).
- Customer uses full-text search and faceted filtering (powered by Product Search Service / OpenSearch), which supports typo-tolerance and category filtering.

**3. Shopping Cart**
- Customer adds books to their cart.
- Cart data is persisted in **Apache Cassandra** by Cart Service, providing a "persistent cart" that survives device switches and session expiry.

**4. Checkout (Order Creation)**
- Customer initiates checkout.
- **Order Service** makes a **synchronous HTTP GET** to Product Service to validate real-time price and stock availability — this is one of the few permitted synchronous inter-service calls because the result is immediately critical.
- If stock is available: a Redis key is set as a **Soft-Booking lock** (TTL: 15 minutes), and an order is created with status `PENDING`.
- Cart data in Cassandra is cleared (soft-deleted) upon successful checkout.

**5. Payment (Xendit Integration)**
- Order Service instructs Payment Service to generate an invoice.
- Payment Service calls the Xendit API to create an Invoice / Virtual Account.
- Customer pays through the Xendit sandbox dashboard.
- Xendit fires a **Webhook (Callback)** back to Payment Service.
- Payment Service validates the webhook signature, then **publishes** an `OrderPaid` event to RabbitMQ (does NOT call Order Service directly).
- Order Service **consumes** the event and updates the order status to `PAID`.

**6. Notification**
- After the `OrderPaid` event is consumed, a second consumer (Notification Service) picks it up and sends a confirmation email.
- The event payload already contains all necessary data (customer email, order details) — Notification Service does **not** call back to Order Service (*Event-Carried State Transfer* pattern).

---

### 3.2 Admin Flow

Admin endpoints are protected by **Coarse-Grained role checking** at the API Gateway level (requests to `/api/admin/*` are blocked unless the JWT `role` is `admin`).

**1. Product Catalog Management**
- Admin can **Create**, **Update** (price, stock), and **Delete** products.
- *Architecture Note:* Every mutation publishes a `ProductUpdated` event to RabbitMQ, which asynchronously syncs the change to OpenSearch.

**2. Order Monitoring**
- Admin can view all orders across all customers.
- Admin can manually advance order status: `PAID` → `SHIPPED` → `DELIVERED`.

---

### 3.3 Automated System Flows (Background Jobs)

The system operates autonomously in the background to maintain data integrity.

**1. Order Expiry**
- A customer has 15 minutes to pay after checkout.
- The Redis Soft-Booking key expires automatically after 15 minutes (TTL), releasing the stock reservation.
- To update the database status from `PENDING` to `EXPIRED`, one of two mechanisms is used:
  - *Primary:* Payment Service listens for the `Invoice Expired` Webhook from Xendit and publishes an `OrderExpired` event to RabbitMQ.
  - *Fallback:* Order Service runs a **Cron Job / Scheduler** every 5 minutes to sweep and expire overdue pending orders.

**2. Search Index Sync (CQRS)**
- On every catalog mutation by Admin, Product Service publishes a `ProductUpdated` event.
- Product Search Service consumes it asynchronously and updates the OpenSearch index, ensuring search results always reflect current pricing.

---

## 4. Scope Constraints

To maintain focus on infrastructure over business logic, the following simplifications are in place:

| Area | Decision |
|---|---|
| **Products** | Books only. Schema: `id`, `title`, `author`, `publisher`, `isbn`, `price`, `stock`, plus `categories` and `book_categories` relational tables. |
| **Shipping** | No location-based shipping cost. Flat rate or free shipping assumed. |
| **Payment** | Xendit Sandbox only. Focus is on the Invoice creation cycle and secure Webhook handling. |
| **Images** | No file uploads. Dummy image URL strings only. |

---

## 5. Security & Protection

All enforcement is at the **API Gateway** level. Internal services are not aware of external threats.

| Layer | Mechanism | Detail |
|---|---|---|
| **IP Rate Limiting** | Redis counter per IP | Max 100 req/min on public endpoints (e.g., `/login`, `/search`) |
| **User Rate Limiting** | Redis counter per User ID | Max 1 req/10s on `/checkout` to eliminate double-click bugs |
| **Blocked Response** | HTTP `429 Too Many Requests` | Request is dropped at the gate; never reaches internal services |
| **Network Isolation** | Docker Private Network | Only API Gateway exposes a public port. All other services are unreachable from outside. |
| **Authorization** | Hybrid Permission Model | Coarse-grained (role check) at Gateway; Fine-grained (data ownership check) inside each service |
