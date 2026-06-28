# Order Service — Service Documentation

**Language:** PHP (Laravel)  
**Store:** PostgreSQL (`order_db`) + Redis  
**Internal Port:** `3005`  
**Owned by:** Commerce Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

This is the **transactional core** of the system — the most complex service in the architecture. It:

- Validates stock and price synchronously at checkout (HTTP call to Product Service)
- Creates a stock reservation lock in Redis (Soft-Booking, TTL: 15 min)
- Persists orders with lifecycle state management
- Consumes `OrderPaid` event from RabbitMQ to update order status
- Runs a background Cron Job to sweep expired pending orders

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/checkout` | ✅ `X-User-ID` | Full checkout flow (validate → lock → create order) |
| `GET` | `/orders` | ✅ `X-User-ID` | List the authenticated user's own orders |
| `GET` | `/admin/orders` | ✅ Admin | List all orders system-wide |
| `PUT` | `/admin/orders/:id/status` | ✅ Admin | Advance order status: `PAID` → `SHIPPED` → `DELIVERED` |

---

## Order Status Lifecycle

```
PENDING → PAID → SHIPPED → DELIVERED
PENDING → EXPIRED  (if payment not received within 15 minutes)
```

---

## Database Schema (`order_db`)

```sql
CREATE TABLE orders (
  id          SERIAL PRIMARY KEY,
  user_id     INTEGER NOT NULL,
  total       INTEGER NOT NULL,
  status      VARCHAR(20) DEFAULT 'PENDING',
  expires_at  TIMESTAMP,
  created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
  id          SERIAL PRIMARY KEY,
  order_id    INTEGER REFERENCES orders(id),
  product_id  INTEGER NOT NULL,
  quantity    INTEGER NOT NULL,
  unit_price  INTEGER NOT NULL             -- snapshot of price at checkout time
);
```

> `unit_price` is stored as a **price snapshot** — the price at the moment of checkout. This prevents price changes from retroactively affecting past orders.

---

## Flow: Checkout & Soft-Booking

```mermaid
sequenceDiagram
    participant C as Client
    participant G as API Gateway
    participant O as Order Service
    participant P as Product Service
    participant R as Redis (Soft-Booking)
    participant DB as order_db (PostgreSQL)

    C->>G: POST /api/checkout { items: [{ product_id: 1, qty: 2 }] }
    G->>O: Forward + X-User-ID: 123
    O->>P: [Sync HTTP GET] /products/1
    P-->>O: { stock: 10, price: 50000 }

    alt Stock insufficient
        O-->>C: 400 Bad Request — out of stock
    else Stock available
        O->>R: SET "soft_book:book_1:user_123" 1 EX 900
        O->>DB: INSERT order (status: PENDING), INSERT order_items
        O->>CartService: DELETE /cart (clear user cart)
        O-->>G: 201 { order_id, expires_at }
        G-->>C: Proceed to payment
    end
```

---

## Flow: Payment Confirmed (Async Consumer)

```mermaid
sequenceDiagram
    participant RMQ as RabbitMQ
    participant O as Order Service
    participant DB as order_db (PostgreSQL)
    participant P as Product Service

    RMQ-)O: consume "OrderPaid" { order_id, ... }
    O->>DB: UPDATE orders SET status = 'PAID' WHERE id = order_id
    O->>P: [Sync HTTP PATCH] /products/1/stock — decrement stock permanently
    O->>DB: DELETE from Redis key "soft_book:book_1:user_123"
```

---

## Background Job: Order Expiry Cron

Runs every **5 minutes** via Laravel Scheduler:

```php
// Pseudo-code
$expiredOrders = Order::where('status', 'PENDING')
                      ->where('expires_at', '<', now())
                      ->get();

foreach ($expiredOrders as $order) {
    $order->update(['status' => 'EXPIRED']);
    // Redis key has already auto-expired via TTL
}
```

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `DB_CONNECTION` | `pgsql` | Laravel DB driver |
| `DATABASE_URL` | `postgres://user:pass@postgres:5432/order_db` | PostgreSQL connection |
| `REDIS_URL` | `redis:6379` | Redis for soft-booking |
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672` | Event consumer |
| `PRODUCT_SERVICE_URL` | `http://product-service:3002` | Sync HTTP target for stock validation |
| `CART_SERVICE_URL` | `http://cart-service:3004` | Cart clearing after checkout |
| `PORT` | `3005` | Internal service port |
