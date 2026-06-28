# API Gateway — Service Documentation

**Language:** Go (Fiber / Gin)  
**Stores:** Redis (rate limit counters)  
**Internal Port:** `8000` (only service with a public-facing port)  
**Owned by:** Infrastructure Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

| # | Responsibility | Detail |
|---|---|---|
| 1 | **Routing** | Forwards public requests to the correct internal service container |
| 2 | **Stateless JWT Validation** | Verifies JWT signature using shared secret key — no call to User Service required |
| 3 | **Header Injection** | Strips JWT, injects `X-User-ID` and `X-User-Role` headers for internal services |
| 4 | **Coarse-Grained Authorization** | Blocks `/api/admin/*` routes if JWT role is not `admin` |
| 5 | **Rate Limiting — IP** | Max 100 req/min per IP address (DDoS / bot protection) |
| 6 | **Rate Limiting — User ID** | Max 1 req/10s per user on transactional endpoints (double-click prevention) |
| 7 | **API Stitching (BFF)** | Fetches data from multiple services and merges into a single response |
| 8 | **Load Test Bypass** | Disables rate limiting if request contains `X-Stress-Test-Bypass` secret header (dev/staging only) |

---

## Routing Table

| Public URL | Forwards To | Auth Required |
|---|---|---|
| `POST /api/auth/login` | User Service | ❌ No |
| `POST /api/auth/register` | User Service | ❌ No |
| `GET /api/products` | Product Service | ❌ No |
| `GET /api/search` | Product Search Service | ❌ No |
| `GET /api/cart` | Cart Service | ✅ Yes |
| `POST /api/cart` | Cart Service | ✅ Yes |
| `POST /api/checkout` | Order Service | ✅ Yes (rate limited: 1/10s) |
| `GET /api/orders` | Gateway BFF (stitches Order + Product) | ✅ Yes |
| `POST /api/admin/*` | Various (Admin role only) | ✅ Yes + Admin |

---

## Flow: Protected Endpoint (e.g., Checkout)

```mermaid
sequenceDiagram
    participant C as Client (App/Web)
    participant G as API Gateway
    participant R as Redis (Rate Limit)
    participant O as Order Service

    C->>G: POST /api/checkout (Bearer JWT)
    G->>R: Increment & check counter for IP + User ID
    alt Limit exceeded
        R-->>G: Blocked
        G-->>C: 429 Too Many Requests
    else Within limit
        G->>G: Verify JWT signature (stateless)
        alt JWT invalid or expired
            G-->>C: 401 Unauthorized
        else JWT valid
            G->>G: Extract user_id and role from JWT payload
            G->>O: Forward request + X-User-ID, X-User-Role headers
            O-->>G: JSON response
            G-->>C: 200 OK
        end
    end
```

---

## Flow: API Stitching — Order History Page

```mermaid
sequenceDiagram
    participant C as Client
    participant G as API Gateway
    participant O as Order Service
    participant P as Product Service

    C->>G: GET /api/orders (Bearer JWT)
    G->>O: GET /orders?user_id=123
    O-->>G: [{ order_id, product_id: 55, amount }, ...]
    G->>P: GET /products?ids=55
    P-->>G: [{ id: 55, title: "Bumi Manusia", cover_url }]
    G->>G: Merge order + product data
    G-->>C: [{ order_id, title: "Bumi Manusia", amount }]
```

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `JWT_SECRET` | `supersecret` | Shared key for JWT signature validation |
| `REDIS_URL` | `redis:6379` | Redis connection for rate limit counters |
| `USER_SERVICE_URL` | `http://user-service:3001` | Internal routing target |
| `ORDER_SERVICE_URL` | `http://order-service:3005` | Internal routing target |
| `STRESS_TEST_BYPASS_KEY` | `k6-bypass-secret` | Secret for disabling rate limit during load tests |
