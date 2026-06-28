# Notification Service — Service Documentation

**Language:** Go or Python  
**Store:** None  
**Internal Port:** None (no HTTP server)  
**Owned by:** Commerce Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

This is a **pure event consumer**. It has no HTTP endpoints and does not respond to any requests. Its entire lifecycle is:

1. Connect to RabbitMQ on startup
2. Listen for events indefinitely
3. For each event: render an email template and send via SMTP

---

## Why No HTTP Endpoints?

Notification Service does not serve user-facing requests. It exists solely to process background jobs triggered by system events. Adding an HTTP server would add unnecessary complexity and resource overhead.

---

## RabbitMQ: Consumed Events

| Event | Action |
|---|---|
| `OrderPaid` | Send order confirmation email to the customer |

> The `OrderPaid` event payload already contains all the data needed (email address, order details, item list). This service **never calls back to Order Service** to fetch additional data — *Event-Carried State Transfer*. See [ADR-004](../adr/004-rabbitmq-async-communication.md).

---

## Email Template: Order Confirmation

Triggered by: `OrderPaid` event

```
Subject: Your order is confirmed — {{ items[0].title }} and more

Hi {{ customer_name }},

Thank you for your purchase! Here's your order summary:

{% for item in items %}
  - {{ item.title }} x{{ item.qty }}  →  Rp{{ item.price | format_currency }}
{% endfor %}

Total: Rp{{ total_amount | format_currency }}

Your books are being prepared for shipment.

— The Bookstore Team
```

---

## Flow: Event Consumed → Email Sent

```mermaid
sequenceDiagram
    participant RMQ as RabbitMQ
    participant N as Notification Service
    participant SMTP as SMTP Server

    RMQ-)N: consume "OrderPaid" { order_id, customer_email, items[] }
    N->>N: Render HTML email template with payload data
    N->>SMTP: Send email to customer_email
    alt SMTP success
        SMTP-->>N: 250 OK
        N->>N: Log { level: "info", message: "Email sent", correlation_id }
    else SMTP failure
        SMTP-->>N: Error
        N->>N: Log { level: "error", message: "Email failed", correlation_id }
        N->>RMQ: NACK — requeue message for retry
    end
```

---

## Error Handling & Retry

- On SMTP failure: message is **NACKed** and returned to the queue for redelivery.
- After 3 failed attempts: message is routed to the **Dead Letter Queue (DLQ)** for manual inspection.
- All failures must be logged in **JSON Structured Logging** format with `correlation_id` for Loki traceability.

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672` | RabbitMQ consumer connection |
| `SMTP_HOST` | `smtp.mailtrap.io` | SMTP server (use Mailtrap for dev) |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USER` | `your-mailtrap-user` | SMTP credentials |
| `SMTP_PASS` | `your-mailtrap-pass` | SMTP credentials |
| `EMAIL_FROM` | `noreply@bookstore.dev` | Sender address |
