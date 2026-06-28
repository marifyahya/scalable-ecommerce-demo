# Product Service — Service Documentation

**Language:** Go (Gin)  
**Store:** PostgreSQL (`product_db`)  
**Internal Port:** `3002`  
**Owned by:** Catalog Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

This service is the **Source of Truth** for all product (book) data. Any mutation here must be propagated to OpenSearch via an async event. No other service writes to `product_db`.

- Book catalog CRUD (Create, Read, Update, Delete)
- Stock quantity management
- Publishing `ProductCreated` / `ProductUpdated` / `ProductDeleted` events to RabbitMQ for CQRS sync

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/products` | ❌ No | List books with pagination |
| `GET` | `/products/:id` | ❌ No | Single book detail (price, stock, categories) |
| `POST` | `/admin/products` | ✅ Admin | Create new book, publish `ProductCreated` event |
| `PUT` | `/admin/products/:id` | ✅ Admin | Update price/stock, publish `ProductUpdated` event |
| `DELETE` | `/admin/products/:id` | ✅ Admin | Soft-delete book, publish `ProductDeleted` event |

---

## Database Schema (`product_db`)

```sql
CREATE TABLE products (
  id          SERIAL PRIMARY KEY,
  title       VARCHAR(255) NOT NULL,
  author      VARCHAR(150) NOT NULL,
  publisher   VARCHAR(150),
  isbn        VARCHAR(20) UNIQUE,
  price       INTEGER NOT NULL,           -- in IDR (Rupiah)
  stock       INTEGER NOT NULL DEFAULT 0,
  cover_url   TEXT,
  deleted_at  TIMESTAMP                   -- soft delete
  created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE categories (
  id    SERIAL PRIMARY KEY,
  name  VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE book_categories (
  book_id     INTEGER REFERENCES products(id),
  category_id INTEGER REFERENCES categories(id),
  PRIMARY KEY (book_id, category_id)
);
```

---

## RabbitMQ Event Payloads

**`ProductCreated` / `ProductUpdated`:**
```json
{
  "event": "ProductUpdated",
  "id": 99,
  "title": "Bumi Manusia",
  "author": "Pramoedya Ananta Toer",
  "price": 85000,
  "stock": 50,
  "categories": ["Fiction", "Historical"]
}
```

> This payload is consumed by **Product Search Service** to update the OpenSearch index. See [product-search-service.md](./product-search-service.md).

---

## Flow: Admin Creates Product → CQRS Sync

```mermaid
sequenceDiagram
    actor A as Admin
    participant G as API Gateway
    participant P as Product Service
    participant DB as product_db (PostgreSQL)
    participant RMQ as RabbitMQ
    participant PS as Product Search Service
    participant OS as OpenSearch

    A->>G: POST /api/admin/products { title, price, stock }
    G->>P: Forward (X-User-Role: admin verified)
    P->>DB: INSERT INTO products
    DB-->>P: { id: 99 }
    P->>RMQ: publish "ProductCreated" { id: 99, title, price, categories }
    P-->>G: 201 Created
    G-->>A: Product saved

    Note over RMQ,OS: Async — independent of the above response
    RMQ-)PS: consume "ProductCreated"
    PS->>OS: PUT /books/_doc/99 { ... }
    OS-->>PS: Indexed
```

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://user:pass@postgres:5432/product_db` | PostgreSQL connection |
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672` | RabbitMQ connection |
| `PORT` | `3002` | Internal service port |
