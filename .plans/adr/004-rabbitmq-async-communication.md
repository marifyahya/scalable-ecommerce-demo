# ADR-004: RabbitMQ for Asynchronous Communication

**Status:** Accepted  
**Date:** 2026-06-28  
**Deciders:** Project Architect

---

## Context

Multiple operations in this system involve cross-service state changes (payment success → order status update → email notification). We considered two approaches:
1. **Synchronous chain:** Service A calls Service B via HTTP, which calls Service C, etc.
2. **Asynchronous events:** Service A publishes an event to a message broker. B and C consume independently.

## Decision

Use **RabbitMQ** as the message broker for all asynchronous cross-service communication.

## Golden Rule Applied

> Asynchronous (RabbitMQ) is mandatory for any write/status-change operation that affects another service's data, or any operation that can be deferred without impacting the user's current session.

## Pattern: Event-Carried State Transfer

Events published to RabbitMQ **must carry all data** the consumers need. A consumer must never need to call back to the event publisher to fetch additional information. This ensures complete consumer independence.

**Example — `OrderPaid` event payload:**
```json
{
  "event": "OrderPaid",
  "order_id": "ORD-999",
  "customer_email": "user@example.com",
  "customer_name": "Budi Santoso",
  "total_amount": 250000,
  "items": [{ "title": "Bumi Manusia", "qty": 1, "price": 250000 }]
}
```

## Why RabbitMQ over Kafka?

| Factor | RabbitMQ | Kafka |
|---|---|---|
| **Complexity** | Low — straightforward setup and management | High — requires Zookeeper/KRaft, topic partitioning strategy |
| **Learning curve** | Manageable for a learning project | Steep; adds operational overhead |
| **Use case fit** | Task queues, event routing, pub/sub | High-throughput log streaming, event sourcing |
| **Verdict** | ✅ Fits our learning scope | Overkill for this domain |

## Consequences

**Positive:**
- Services are fully decoupled. If Notification Service crashes, the `OrderPaid` message is retained in RabbitMQ and redelivered upon recovery — no emails are permanently lost.
- RabbitMQ acts as a traffic buffer during flash sales, absorbing thousands of checkout events without overwhelming Order Service.

**Negative / Trade-offs:**
- **Eventual consistency:** There is a small lag between payment confirmation and the order status update visible in the UI.
- Requires a **Dead Letter Queue (DLQ)** strategy for messages that fail processing repeatedly.
- Adds RabbitMQ as a dependency that must be operational for the async flows to function.
