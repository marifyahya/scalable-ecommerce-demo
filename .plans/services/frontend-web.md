# Frontend Web — Service Documentation

**Language:** TypeScript (Svelte + Vite)  
**Store:** Browser Storage (Cookies for JWT, LocalStorage for UI State)  
**Internal Port:** `3000` (Dev Server) / Static files (Prod)  
**Owned by:** Web Experience Team

> For cross-service communication rules and the full system diagram, see [blueprint.md](../blueprint.md).

---

## Responsibilities

This is a **Single Page Application (SPA)** that serves as the face of our e-commerce platform. It does not possess any backend logic or database connections of its own.

- Renders the UI using Svelte components.
- Manages client-side routing.
- Makes HTTP REST calls exclusively to the **API Gateway** (it never calls User Service or Order Service directly).
- Handles responsive styling via Tailwind CSS.

---

## Why Svelte?

Unlike React or Vue, Svelte is a compiler. It does not ship a bulky Virtual DOM runtime to the browser. Instead, it compiles `.svelte` components into highly optimized, surgically precise vanilla JavaScript.

This makes Svelte the absolute lightest and fastest modern framework available, aligning perfectly with our architecture's philosophy of minimizing overhead. 

---

## Interaction with API Gateway

The Frontend assumes the API Gateway acts as a **Backend-for-Frontend (BFF)**.

**1. Authentication:**
When a user logs in, the Gateway returns a JWT. The Frontend must store this JWT securely (preferably in an `httpOnly` cookie or memory) and attach it as a `Bearer` token to all subsequent authenticated requests.

**2. Data Stitching:**
The Frontend does not need to fetch data from multiple services and combine them. For example, to render the "Order History" page, the Frontend simply calls `GET /api/orders` on the Gateway, and the Gateway does the heavy lifting of stitching Order records with Product details.

---

## Core Pages / Routes

| Path | Description | API Gateway Calls |
|---|---|---|
| `/` | Homepage (Product Catalog) | `GET /api/products` |
| `/search` | Advanced Search | `GET /api/search` |
| `/login` | User Login | `POST /api/auth/login` |
| `/cart` | Shopping Cart | `GET /api/cart`, `POST /api/cart` |
| `/checkout` | Checkout & Payment | `POST /api/checkout` |
| `/orders` | Order History | `GET /api/orders` |

---

## Environment Variables

| Variable | Example | Description |
|---|---|---|
| `VITE_API_GATEWAY_URL` | `http://localhost:8000` | The public URL of the API Gateway |
