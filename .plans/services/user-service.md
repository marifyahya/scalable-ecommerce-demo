# User Service — Service Documentation

**Language:** Node.js (NestJS + TypeScript)  
**Store:** PostgreSQL (`user_db`)  
**Internal Port:** `3001`  
**Owned by:** Identity Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

This service is the **sole authority for identity and JWT issuance** in the system. No other service may generate JWT tokens.

- User registration and login
- Password hashing (`bcrypt`) and verification
- JWT token generation and signing
- User profile management

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/register` | ❌ No | Hash password, persist user, return success |
| `POST` | `/login` | ❌ No | Verify password, **generate and sign JWT token** |
| `GET` | `/users/me` | ✅ `X-User-ID` header | Return current user profile |

---

## Database Schema (`user_db`)

```sql
CREATE TABLE users (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  email       VARCHAR(150) UNIQUE NOT NULL,
  password    VARCHAR(255) NOT NULL,   -- bcrypt hash
  role        VARCHAR(20) DEFAULT 'customer',
  created_at  TIMESTAMP DEFAULT NOW()
);
```

---

## JWT Token Structure

```json
{
  "user_id": 123,
  "role": "customer",
  "iat": 1719000000,
  "exp": 1719003600
}
```

> Token expiry: **1 hour**. The secret key is shared with API Gateway for stateless validation. See [ADR-005](../adr/005-stateless-jwt-at-gateway.md).

---

## Flow: Login & JWT Issuance

```mermaid
sequenceDiagram
    participant C as Client
    participant G as API Gateway
    participant U as User Service
    participant DB as user_db (PostgreSQL)

    C->>G: POST /api/auth/login { email, password }
    G->>U: Forward (unauthenticated route — no JWT check)
    U->>DB: SELECT * FROM users WHERE email = ?
    DB-->>U: { id, email, password_hash, role }
    U->>U: bcrypt.compare(input_password, password_hash)
    alt Password mismatch
        U-->>G: 401 Invalid credentials
        G-->>C: Login failed
    else Password correct
        U->>U: jwt.sign({ user_id, role }, JWT_SECRET, { expiresIn: '1h' })
        U-->>G: 200 { token: "eyJhbG..." }
        G-->>C: Login successful — client stores token
    end
```

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://user:pass@postgres:5432/user_db` | PostgreSQL connection string |
| `JWT_SECRET` | `supersecret` | Signing key — must match API Gateway's value |
| `BCRYPT_ROUNDS` | `12` | bcrypt salt rounds |
| `PORT` | `3001` | Internal service port |
