# ADR-005: Stateless JWT Validation at API Gateway

**Status:** Accepted  
**Date:** 2026-06-28  
**Deciders:** Project Architect

---

## Context

Every authenticated request needs to be verified. We had two options:

1. **Stateful validation:** API Gateway forwards the JWT to User Service on every request. User Service queries the database to confirm validity.
2. **Stateless validation:** API Gateway validates the JWT signature independently using a shared secret key, without calling any other service.

## Decision

Use **Stateless JWT Validation** at the API Gateway using a shared secret key (HS256) or a public/private key pair (RS256).

## Rationale

Option 1 (stateful) makes User Service a **synchronous bottleneck** on every single request in the system. If User Service goes down or slows down, authentication fails across all services — a catastrophic single point of failure.

With Option 2 (stateless), the Gateway validates the JWT signature locally in microseconds. User Service is only called for:
- Initial login (token issuance)
- Profile updates

## Token Issuance Rule

**User Service is the sole authority for issuing JWT tokens.** The Gateway and User Service share the same secret key. The Gateway can verify tokens but cannot issue new ones.

## Header Injection (Internal Auth)

Upon successful JWT validation, the Gateway:
1. Strips the `Authorization: Bearer <jwt>` header from the forwarded request.
2. Injects plain `X-User-ID` and `X-User-Role` headers.
3. Internal services trust these headers implicitly — they are unreachable from the public internet (network isolation ensures this).

## Consequences

**Positive:**
- Authentication adds near-zero latency overhead (cryptographic verify ≈ microseconds).
- User Service is never a bottleneck for normal request processing.
- Internal services are completely free from JWT parsing logic.

**Negative / Trade-offs:**
- **Token revocation is not immediate.** A compromised JWT remains valid until its `exp` (expiry) timestamp. Mitigation: use short expiry times (e.g., 15 minutes) with refresh token rotation.
- All services that need to issue or validate tokens must share the same secret — key rotation requires coordinated redeployment.
