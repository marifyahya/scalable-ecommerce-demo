# Product Search Service — Service Documentation

**Language:** Python (FastAPI)  
**Store:** OpenSearch  
**Internal Port:** `3003`  
**Owned by:** Catalog Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

This service is a **read-optimized replica** of the product catalog, purpose-built for search. It does **not** write to `product_db`. It only:

- Receives sync events from Product Service via RabbitMQ (consumer)
- Indexes product documents into OpenSearch
- Serves search queries (full-text, typo-tolerant, faceted)

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/search` | ❌ No | Full-text search with optional facets |

**Example query:**
```
GET /search?q=pramoedya&category=Fiction&min_price=50000&max_price=150000
```

**Example response:**
```json
{
  "total": 3,
  "results": [
    {
      "id": 99,
      "title": "Bumi Manusia",
      "author": "Pramoedya Ananta Toer",
      "price": 85000,
      "categories": ["Fiction", "Historical"]
    }
  ]
}
```

---

## RabbitMQ: Consumed Events

| Event | Action |
|---|---|
| `ProductCreated` | Insert new document into OpenSearch |
| `ProductUpdated` | Update existing document in OpenSearch |
| `ProductDeleted` | Remove document from OpenSearch index |

---

## OpenSearch Index Mapping

The relational `book_categories` table is **flattened** into an array for OpenSearch:

```json
{
  "mappings": {
    "properties": {
      "id":         { "type": "integer" },
      "title":      { "type": "text" },
      "author":     { "type": "text" },
      "price":      { "type": "integer" },
      "stock":      { "type": "integer" },
      "categories": { "type": "keyword" }
    }
  }
}
```

---

## Flow: Product Sync + Search Query

```mermaid
sequenceDiagram
    participant RMQ as RabbitMQ
    participant PS as Product Search Service
    participant OS as OpenSearch
    participant C as Client

    Note over RMQ,OS: Background sync (always running)
    RMQ-)PS: consume "ProductUpdated" { id: 99, title, price, categories }
    PS->>PS: Flatten categories array
    PS->>OS: PUT /books/_doc/99 { title, price, categories[] }
    OS-->>PS: Acknowledged

    Note over C,OS: User performs a search
    C->>PS: GET /search?q=pramoedya&category=Fiction
    PS->>OS: POST /books/_search { query: { multi_match } }
    OS-->>PS: [{ id: 99, title, price, ... }]
    PS-->>C: 200 { total: 3, results: [...] }
```

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `OPENSEARCH_URL` | `http://opensearch:9200` | OpenSearch node URL |
| `RABBITMQ_URL` | `amqp://guest:guest@rabbitmq:5672` | RabbitMQ connection |
| `OPENSEARCH_INDEX` | `books` | Index name |
| `PORT` | `3003` | Internal service port |
