# ЕщёЕсть API

Go modular-monolith API for the pet-project backend. It currently provides:

- PostgreSQL migrations and seed offers;
- six-digit email verification for registration and passwordless login;
- hashed challenges and server-side, revocable HttpOnly sessions;
- current-user read/update endpoints;
- public offer reads and role-protected offer CRUD;
- JSON logs, request IDs, input limits, and health checks.

The authoritative initial contract is [`api/openapi.yaml`](api/openapi.yaml).

## Local run

Go 1.26+ and PostgreSQL 14+ are supported. Set the environment and start the API:

```bash
export DATABASE_URL='postgres://user:password@127.0.0.1:5432/eshche_est?sslmode=disable'
export OTP_PEPPER='replace-with-at-least-32-random-characters'
export EMAIL_DELIVERY='development'
export EXPOSE_DEV_CODES='true'
go run ./cmd/api
```

Development codes are returned only when both development delivery and
`EXPOSE_DEV_CODES=true` are configured. SMTP mode requires port 587 with
STARTTLS and never returns or logs verification codes.
