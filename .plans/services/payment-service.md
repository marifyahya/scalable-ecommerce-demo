# Payment Service — Service Documentation

**Language:** Node.js (Express)  
**Store:** PostgreSQL (`payment_db`)  
**Internal Port:** `3006` (+ public webhook receiver port)  
**Owned by:** Payments Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

- Generate payment invoices via the Xendit API
- Receive and validate Xendit webhook callbacks
- Publish `OrderPaid` or `OrderExpired` events to RabbitMQ
- **Does NOT call Order Service directly** — all communication is via events

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/payments` | ✅ Internal (from Order Service) | Call Xendit API to generate invoice, return payment URL |
| `POST` | `/xendit/callback` | ❌ Public (Xendit only) | Receive webhook, validate, publish event to RabbitMQ |

---

## Database Schema (`payment_db`)

```sql
CREATE TABLE payment_records (
  id              SERIAL PRIMARY KEY,
  order_id        INTEGER NOT NULL,
  xendit_id       VARCHAR(100) UNIQUE,      -- Xendit invoice ID
  amount          INTEGER NOT NULL,
  status          VARCHAR(20) DEFAULT 'PENDING',
  payment_method  VARCHAR(50),
  paid_at         TIMESTAMP,
  created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## Xendit Webhook Validation

Xendit sends a signature in the `x-callback-token` header. Payment Service must validate this before processing:

```js
// Pseudo-code
const incomingToken = req.headers['x-callback-token'];
if (incomingToken !== process.env.XENDIT_WEBHOOK_TOKEN) {
  return res.status(401).json({ error: 'Invalid webhook token' });
}
```

---

## RabbitMQ Events Published

| Event | Trigger | Payload |
|---|---|---|
| `OrderPaid` | Xendit webhook `status: PAID` | `{ order_id, user_email, amount, items[] }` |
| `OrderExpired` | Xendit webhook `status: EXPIRED` | `{ order_id }` |

**`OrderPaid` Payload (Event-Carried State Transfer):**
```json
{
  "event": "OrderPaid",
  "order_id": "ORD-999",
  "customer_email": "user@example.com",
  "customer_name": "Budi Santoso",
  "total_amount": 250000,
  "items": [
    { "title": "Bumi Manusia", "qty": 1, "price": 250000 }
  ]
}
```

> Payload carries **all data** consumers need. Notification Service does not need to call back to Order Service. See [ADR-004](../adr/004-rabbitmq-async-communication.md).

---

## Flow: Webhook Received → Event Published

```mermaid
sequenceDiagram
    participant X as Xendit
    participant P as Payment Service
    participant DB as payment_db (PostgreSQL)
    participant RMQ as RabbitMQ

    X->>P: POST /xendit/callback { invoice_id, status: "PAID", ... }
    P->>P: Validate x-callback-token header
    alt Invalid token
        P-->>X: 401 Unauthorized
    else Valid
        P->>DB: UPDATE payment_records SET status = 'PAID', paid_at = NOW()
        P->>P: Build OrderPaid event payload (enrich with order items)
        P->>RMQ: publish "OrderPaid" { order_id, email, amount, items[] }
        P-->>X: 200 OK
    end
```

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://user:pass@postgres:5432/payment_db` | PostgreSQL connection |
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672` | RabbitMQ for event publishing |
| `XENDIT_SECRET_KEY` | `xnd_development_...` | Xendit API key (sandbox) |
| `XENDIT_WEBHOOK_TOKEN` | `your-callback-token` | Validation token for incoming webhooks |
| `ORDER_SERVICE_URL` | `http://order-service:3005` | Used only to fetch order details for event payload |
| `PORT` | `3006` | Internal service port |
