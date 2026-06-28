# ADR-006: Soft-Booking via Redis for Stock Reservation

**Status:** Accepted  
**Date:** 2026-06-28  
**Deciders:** Project Architect

---

## Context

During checkout, there is a time gap between when a user confirms their cart and when their payment clears. Without a reservation mechanism, two users could simultaneously checkout the last copy of a book — both receive payment instructions, but only one can actually be fulfilled. This is the **overselling problem**.

We considered two approaches:
1. **Hard deduction:** Permanently reduce stock in PostgreSQL at checkout, restore it if payment fails.
2. **Soft-booking:** Reserve stock temporarily (without permanent deduction) using a TTL-based lock.

## Decision

Use **Redis with TTL** as a soft-booking (stock reservation) mechanism.

A Redis key is set at checkout time with a 15-minute expiration:
```
SET "soft_book:{book_id}:user_{user_id}" 1 EX 900
```

## Rationale

| Factor | Hard Deduction | Soft-Booking via Redis |
|---|---|---|
| **Rollback on failure** | Requires compensating transaction in PostgreSQL | Key simply expires — no cleanup needed |
| **Speed** | PostgreSQL UPDATE + potential retry | Redis SET is a sub-millisecond in-memory operation |
| **Oversell protection** | ✅ Yes | ✅ Yes |
| **Complexity** | Higher (requires saga/rollback logic) | Lower (Redis TTL handles expiry automatically) |

## Flow

1. **Checkout:** `SET "soft_book:book_1:user_123" EX 900` → Order saved as `PENDING`
2. **Payment success (within 15 min):** Webhook received → Stock permanently deducted in PostgreSQL → Redis key deleted
3. **Payment timeout:** Redis TTL expires → Stock automatically released → Order status set to `EXPIRED` by Cron Job or Xendit webhook

## Consequences

**Positive:**
- Zero rollback logic required for the expiry case — Redis handles it automatically.
- Prevents overselling with near-zero performance cost.
- The Cron Job fallback in Order Service (`laravel schedule:run`) ensures database cleanup even if the Xendit webhook is delayed.

**Negative / Trade-offs:**
- If Redis goes down between checkout and payment confirmation, the soft-booking is lost. Stock may become available to others during that window.
- There is no hard database lock — under extreme race conditions (millions of concurrent checkouts), multiple users could pass the Redis check before any key is set. Mitigation: use Redis atomic `SET NX` (set if not exists) with appropriate granularity.
