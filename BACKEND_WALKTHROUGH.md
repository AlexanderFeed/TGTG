# Backend walkthrough for a Go beginner

This document explains how the ЕщёЕсть frontend, Go API, and PostgreSQL database
communicate. Read it beside the code. For commands to start the project locally
or against the VPS database, use [STARTUP.md](STARTUP.md).

## Start here

The shortest useful reading order is:

1. [`apps/api/cmd/api/main.go`](apps/api/cmd/api/main.go) — backend program entry point.
2. [`apps/api/internal/httpapi/server.go`](apps/api/internal/httpapi/server.go) — middleware and every Go route.
3. [`apps/web/app/composables/useAuth.ts`](apps/web/app/composables/useAuth.ts) — frontend auth calls.
4. [`apps/web/server/api/[...path].ts`](apps/web/server/api/[...path].ts) — Nuxt-to-Go proxy.
5. [`apps/api/internal/httpapi/auth.go`](apps/api/internal/httpapi/auth.go) — HTTP auth workflow.
6. [`apps/api/internal/store/store.go`](apps/api/internal/store/store.go) — auth/profile SQL and transaction.
7. [`apps/api/migrations/00001_initial.sql`](apps/api/migrations/00001_initial.sql) — database tables.

Do not try to understand every line on the first pass. First follow one request
from the browser to PostgreSQL and back. Registration is the best example.

## The whole application in one picture

```text
Browser
  |
  | http://localhost:3000/register
  v
Nuxt/Vue page: app/pages/register.vue
  |
  | calls requestRegistrationCode(...)
  v
Nuxt composable: app/composables/useAuth.ts
  |
  | POST /api/v1/auth/register/request
  v
Nuxt server proxy: server/api/[...path].ts
  |
  | POST http://api:8080/v1/auth/register/request
  v
Go router: internal/httpapi/server.go
  |
  | calls server.requestRegistration
  v
Go handler: internal/httpapi/auth.go
  |                         |
  | store.UserExists(...)   | mailer.SendVerificationCode(...)
  v                         v
Go store: internal/store    Development sender or SMTP sender
  |
  | parameterized SQL
  v
PostgreSQL
```

The browser never connects to PostgreSQL. Only Go knows the database URL and
database password.

## What “files communicate” means in Go

Files do not send messages to one another automatically. Communication happens
through packages, imported names, function calls, interfaces, and returned
values.

All `.go` files in the same directory normally declare the same package. The Go
compiler treats those files as one package:

```text
internal/httpapi/server.go   package httpapi
internal/httpapi/auth.go     package httpapi
internal/httpapi/offers.go   package httpapi
internal/httpapi/response.go package httpapi
```

That is why `auth.go` can call `writeJSON` from `response.go` without importing
the file. Both functions belong to package `httpapi`.

Code in a different package must import the package and use an exported name:

```go
import "github.com/eshche-est/eshche-est/apps/api/internal/store"

dataStore := store.New(db)
```

In Go, a name beginning with a capital letter is exported from its package.
`store.New` is visible to other packages; `store.normalizeSearch` is private to
package `store`.

### Current package responsibilities

| Package/directory | Responsibility | Calls or is called by |
|---|---|---|
| `cmd/api` | Starts and stops the program; connects dependencies | Calls all setup packages |
| `internal/config` | Reads and validates environment variables | Called by `main` |
| `internal/database` | Opens the PostgreSQL pool and runs migrations | Called by `main` |
| `internal/httpapi` | Routes, middleware, JSON, validation, access rules | Called by `net/http`; calls store/mailer |
| `internal/store` | Parameterized SQL and transactions | Called by handlers; calls PostgreSQL |
| `internal/mailer` | Development or SMTP email delivery | Called by auth handlers |
| `migrations` | Versioned SQL schema embedded in the binary | Used by `database.Migrate` |

This direction is intentional:

```text
HTTP handler -> store -> PostgreSQL
HTTP handler -> mailer -> SMTP
```

The store does not import `httpapi`, and PostgreSQL code does not know anything
about Vue forms. This prevents circular dependencies and keeps each layer small.

## Backend entry point

Go executable programs require:

```go
package main

func main() {
    // Program starts here.
}
```

The backend entry point is [`apps/api/cmd/api/main.go`](apps/api/cmd/api/main.go).
When Docker starts the API, the compiled program runs `main()`.

`main()` performs these steps in order:

1. Creates the JSON logger.
2. Calls `config.Load()` to read environment variables.
3. Creates a root cancellation context for Ctrl+C/Docker shutdown.
4. Calls `database.Open()` to create and test a PostgreSQL pool.
5. Calls `database.Migrate()` to create/update tables.
6. Chooses `DevelopmentSender` or `SMTPSender`.
7. Creates `store.Store`, which holds the database pool.
8. Calls `httpapi.New(...)`, which returns the router.
9. Starts `http.Server` on port `8080` by default.
10. Waits for a shutdown signal and shuts down gracefully.

The important construction line is conceptually:

```go
handler := httpapi.New(cfg, store.New(db), sender, log)
```

This is dependency injection without a framework. `main` creates real
dependencies and passes them to the code that needs them.

## Frontend entry points into the backend

There is not one frontend entry point. Each user action begins in a Vue page and
calls a function from `useAuth()`.

| User action | Vue function | Browser request | Go handler |
|---|---|---|---|
| Request register code | `register.vue -> requestCode` | `POST /api/v1/auth/register/request` | `requestRegistration` |
| Verify registration | `register.vue -> verifyCode` | `POST /api/v1/auth/register/verify` | `verifyRegistration` |
| Request login code | `login.vue -> requestCode` | `POST /api/v1/auth/login/request` | `requestLogin` |
| Verify login | `login.vue -> verifyCode` | `POST /api/v1/auth/login/verify` | `verifyLogin` |
| Load current session | `useAuth -> refreshSession` | `GET /api/v1/auth/me` | `me` after `requireAuth` |
| Edit profile | `profile.vue -> save` | `PATCH /api/v1/users/me` | `updateProfile` after `requireAuth` |
| Log out | `useAuth -> logout` | `POST /api/v1/auth/logout` | `logout` after `requireAuth` |

### Why browser URLs contain `/api`, but Go routes do not

The browser calls:

```text
/api/v1/auth/me
```

Nuxt recognizes `/api/*` as a server route and runs
[`apps/web/server/api/[...path].ts`](apps/web/server/api/[...path].ts). The
catch-all `path` is `v1/auth/me`, so it forwards to:

```text
http://api:8080/v1/auth/me        # Docker
http://127.0.0.1:8080/v1/auth/me  # manual local run
```

Go registers only `/v1/auth/me` in `httpapi/server.go`.

This proxy provides one public origin. The visitor only needs port 3000/HTTPS,
the internal Go container name stays private, and the session cookie works
without cross-origin browser configuration.

## Registration request, line by line across the layers

### 1. Vue collects input

[`apps/web/app/pages/register.vue`](apps/web/app/pages/register.vue) uses Vue
reactive state:

```ts
const form = reactive({ name: '', email: '' })
```

`v-model="form.email"` keeps the input and this state synchronized. The form has
`@submit.prevent="requestCode"`, so Vue calls `requestCode()` without allowing a
normal browser page reload.

### 2. The composable performs HTTP

The page calls `requestRegistrationCode` from
[`useAuth.ts`](apps/web/app/composables/useAuth.ts). It uses `$fetch`:

```ts
$fetch('/api/v1/auth/register/request', {
  method: 'POST',
  headers: { 'X-Requested-With': 'eshche-est-web' },
  body: { name, email },
})
```

Nuxt serializes `body` as JSON. The custom header is required by Go middleware
for state-changing requests.

### 3. Nuxt proxies the request

`server/api/[...path].ts` forwards the method, body, headers, cookies, query
string, status, response headers, and response body to/from Go.

### 4. Go middleware wraps the handler

In [`server.go`](apps/api/internal/httpapi/server.go), `httpapi.New` creates the
Chi router. Before the registration handler runs, global middleware:

1. Assigns a request ID.
2. Recovers from an unexpected panic.
3. applies a 20-second timeout.
4. Adds response security/cache headers.
5. Limits the request body to 1 MiB.
6. Logs method, path, status, and duration.
7. Requires the mutation header for POST/PUT/PATCH/DELETE.

Then Chi matches:

```go
r.Post("/register/request", server.requestRegistration)
```

### 5. The handler validates JSON

`requestRegistration` in [`auth.go`](apps/api/internal/httpapi/auth.go):

1. Calls `decodeJSON` into `requestCodeInput`.
2. Normalizes/validates email and name.
3. Calls `s.store.UserExists(...)`.
4. Calls the shared `s.requestCode(...)` workflow.

`decodeJSON` rejects unknown JSON fields. If the frontend accidentally sends
`userEmail` instead of `email`, the API returns a visible `400` error instead of
silently ignoring the mistake.

### 6. The store reads/writes PostgreSQL

`UserExists` runs a parameterized query:

```sql
SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)
```

`$1` is a placeholder. The email value is sent separately to the database
driver. Never construct SQL by joining user input into the SQL string.

`CreateChallenge` inserts a row into `email_challenges`, containing:

- challenge UUID;
- normalized email;
- purpose (`register` or `login`);
- pending registration name;
- hash of the six-digit code;
- attempts/maximum attempts;
- expiration time.

It does not store the plaintext code.

### 7. Email or local development code

The handler calls the `mailer.Sender` interface:

```go
type Sender interface {
    SendVerificationCode(ctx context.Context, recipient, code, purpose string) error
}
```

- `DevelopmentSender` does nothing. With `EXPOSE_DEV_CODES=true`, the response
  contains `devCode`, and the frontend displays/autofills it.
- `SMTPSender` connects using STARTTLS and sends a real email. SMTP mode cannot
  expose `devCode`.

The response returns to Vue through the same route in reverse.

## Registration verification and transaction

The second form sends `challengeId`, `email`, and the six-digit `code`.
`verifyCode` in `httpapi/auth.go` generates a random session token and sends only
its hash to `Store.CompleteChallenge`.

`CompleteChallenge` starts a PostgreSQL transaction:

```text
BEGIN
  SELECT challenge ... FOR UPDATE
  compare stored code hash with submitted code hash
  INSERT user
  UPDATE challenge SET consumed_at = now()
  INSERT session
COMMIT
```

Why a transaction matters: imagine `users` was inserted but session insertion
failed. Without a transaction, registration would return an error while leaving
the email registered. Retrying would then say the account already exists. With
the transaction, all required writes succeed together or all are rolled back.

`FOR UPDATE` locks the challenge row during the transaction. If two requests
submit the same one-time code simultaneously, they cannot both consume it.

After commit, Go sends:

```text
Set-Cookie: eshche_est_session=<random-token>; HttpOnly; SameSite=Lax; ...
```

The raw random token exists in the browser cookie. PostgreSQL stores only its
SHA-256 hash. `HttpOnly` means Vue JavaScript cannot read the cookie; the browser
attaches it to later requests automatically.

## How an authenticated request works

Profile loading is a good example:

```text
Browser GET /api/v1/auth/me (cookie attached automatically)
  -> Nuxt proxy forwards cookie
  -> Go requireAuth middleware reads cookie
  -> Go hashes raw token
  -> Store.UserBySession queries live session + user
  -> middleware puts User into request context
  -> me handler reads currentUser(request)
  -> JSON { "user": ... }
  -> useAuth updates shared frontend state
```

`context.Context` has two uses here:

1. Cancellation/deadlines: SQL stops if the HTTP request is cancelled.
2. Request-scoped data: `requireAuth` attaches the authenticated user for the
   next handler.

Do not use context as a global variable store. Values should belong to this one
request and travel down its call chain.

## Profile update flow

[`profile.vue`](apps/web/app/pages/profile.vue) declares:

```ts
definePageMeta({ middleware: 'auth' })
```

Nuxt therefore runs [`app/middleware/auth.ts`](apps/web/app/middleware/auth.ts)
before rendering the profile. If auth state has not been checked, it calls
`refreshSession`; if no user exists, it redirects to `/login`.

Saving calls `useAuth.updateProfile`, which sends only `{name, city}`. Go obtains
the user ID from the authenticated session context, not from editable frontend
JSON. This prevents one user from changing another user merely by submitting a
different ID.

## Offer CRUD and the current frontend boundary

Go already exposes:

| Method | Go path | Access |
|---|---|---|
| `GET` | `/v1/offers` | Public |
| `GET` | `/v1/offers/{offerID}` | Public |
| `POST` | `/v1/offers` | Merchant/admin |
| `PUT` | `/v1/offers/{offerID}` | Owner merchant/admin |
| `DELETE` | `/v1/offers/{offerID}` | Owner merchant/admin |

The important current limitation is that discover//offer Vue pages still
import static mock data from [`apps/web/app/data/marketplace.ts`](apps/web/app/data/marketplace.ts).
[`useMarketplace.ts`](apps/web/app/composables/useMarketplace.ts) also stores
favorite IDs in a browser cookie. Those catalog pages do not call Go yet.

Therefore:

- registration, login, session, logout, and profile changes use PostgreSQL;
- direct offer API calls use PostgreSQL;
- visible catalog cards still use mock TypeScript data.

Connecting the catalog UI to `GET /api/v1/offers` is the next frontend/backend
integration milestone in `PLAN.md`.

## PostgreSQL schema and migrations

[`00001_initial.sql`](apps/api/migrations/00001_initial.sql) contains the first
Goose migration.

```text
users
  ^
  | user_sessions.user_id
  |
user_sessions

email_challenges (temporary auth workflow)

offers.created_by -> users.id (nullable for seed offers)
```

`-- +goose Up` creates the schema. `-- +goose Down` reverses it for development.
Go embeds migration files with `//go:embed *.sql`, so the compiled API contains
the SQL it needs.

When the schema changes, add a new migration such as:

```text
00002_create_orders.sql
```

Do not edit an old migration after it has run on your VPS. Goose records that
version as complete and will not rerun the changed file.

### Where your information is stored

| Action | Table/change |
|---|---|
| Request email code | Insert `email_challenges` |
| Wrong code | Increment `email_challenges.attempts` |
| Register successfully | Insert `users`, consume challenge, insert `user_sessions` |
| Login successfully | Consume challenge, insert `user_sessions` |
| Edit profile | Update `users` |
| Log out | Set `user_sessions.revoked_at` |
| Delete offer | Set `offers.status = 'deleted'` |

Use the read-only inspection commands in [STARTUP.md](STARTUP.md#inspect-stored-data)
to see these records locally.

## Go concepts used in this backend

### Struct

A struct groups fields:

```go
type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}
```

The JSON tags control encoded field names. They do not create database mapping;
the store explicitly scans SQL columns into fields.

### Method and pointer receiver

```go
func (s *Store) UserExists(...) ...
```

This is a method on `*Store`. `s` is like a named receiver and provides access
to `s.db`. A pointer receiver avoids copying and shares the same store/pool.

### Interface

```go
type Sender interface {
    SendVerificationCode(...) error
}
```

Go interfaces are satisfied implicitly. Both development and SMTP sender types
implement that method, so both can be passed as `mailer.Sender`.

### Multiple return values and errors

```go
user, err := s.store.UserBySession(...)
if err != nil {
    // handle failure
}
```

Go normally returns errors as values. `errors.Is(err, store.ErrNotFound)` checks
both the error and wrapped errors in its chain.

### `defer`

`defer db.Close()` schedules cleanup when the surrounding function returns.
`defer tx.Rollback()` is especially useful: every early transaction error path
is safe, while rollback after a successful commit is harmless.

### Goroutine

`go func() { server.ListenAndServe() }()` starts the blocking HTTP listener
concurrently. The main goroutine can then wait for an operating-system signal.
The standard HTTP server itself starts request-handling goroutines.

### `database/sql`

`*sql.DB` is a pool, not one permanent connection. Use `QueryRowContext` for one
row, `QueryContext` plus `rows.Next()` for many rows, and `ExecContext` when no
row data is returned.

## How to trace and debug one request

1. Open browser developer tools, select **Network**, submit the form, and inspect
   `/api/v1/...`: method, JSON request, status, and response.
2. Read API logs:

   ```powershell
   docker compose -f compose.yaml -f compose.local.yaml logs -f api
   ```

3. Use the logged `request_id` to correlate the request with an internal error.
4. Find the route path in `internal/httpapi/server.go`.
5. Follow the handler method to `auth.go` or `offers.go`.
6. Follow every `s.store.SomeMethod` call into `internal/store`.
7. Open the migration to confirm table/column constraints.
8. Inspect rows with the `psql` commands in `STARTUP.md`.

Useful status meaning:

| Status | Meaning in this API |
|---|---|
| `200` | Request succeeded |
| `202` | Email challenge accepted/created |
| `204` | Succeeded with no JSON body, such as logout/delete |
| `400` | Invalid JSON/path input |
| `401` | Login/code/session is invalid |
| `403` | Authenticated/requested but not allowed |
| `409` | Email already registered |
| `422` | Valid JSON shape but invalid field values |
| `429` | Cooldown or attempt limit |
| `500` | Unexpected server/store/mailer failure |
| `503` | Health check cannot reach PostgreSQL |

## How to add a new backend feature

For an `orders` feature, work through the layers in this order:

1. Add a new numbered migration with the `orders` table and constraints.
2. Add safe store model/input types.
3. Add parameterized queries in `internal/store/orders.go`.
4. Add HTTP input validation and handlers in `internal/httpapi/orders.go`.
5. Register routes in `internal/httpapi/server.go`; place protected ones inside
   the `requireAuth` group.
6. Document the JSON contract in `apps/api/api/openapi.yaml`.
7. Add a frontend composable that calls `/api/v1/orders`.
8. Call the composable from the relevant Vue page.
9. Test success, invalid input, unauthenticated access, unauthorized access, and
   database failure/rollback behavior.
10. Update `PLAN.md` after the slice is working.

Keep these boundaries:

- Vue should never receive database credentials or issue SQL.
- HTTP handlers should not contain raw SQL.
- Store methods should not write HTTP responses.
- Never trust role/user ID sent by the frontend.
- Never store plaintext session tokens or verification codes in PostgreSQL.
- Use integer kopecks for money, not floating-point rubles.

## Related project documents

- [PLAN.md](PLAN.md) — roadmap and current technical boundary.
- [STARTUP.md](STARTUP.md) — local Docker, VPS database, deployment, and SQL inspection.
- [apps/api/README.md](apps/api/README.md) — concise API overview.
- [apps/api/api/openapi.yaml](apps/api/api/openapi.yaml) — machine-readable HTTP contract.
- [DEPLOY.md](DEPLOY.md) — deployment details.

