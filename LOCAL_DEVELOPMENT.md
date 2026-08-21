# Local development and database storage

The local stack contains three containers:

- `web` — Nuxt at `http://localhost:3000`;
- `api` — Go API at `http://127.0.0.1:8080`;
- `db` — PostgreSQL at `127.0.0.1:5433`.

PostgreSQL data is stored in the named Docker volume
`eshche-est-local_postgres_local_data`, so it survives container restarts and
ordinary `docker compose down` operations.

## Start everything

From the repository root:

```powershell
docker compose -f compose.yaml -f compose.local.yaml up -d --build
docker compose -f compose.yaml -f compose.local.yaml ps
```

Open `http://localhost:3000/register`. Use any valid email address. Local email
delivery is intentionally disabled: the generated six-digit code is returned
to the frontend, displayed, and filled into the code field automatically.

After registering, test logout and passwordless login at
`http://localhost:3000/login`. Use the same email. The local admin account is:

```text
admin@example.test
```

Registering that address grants the `admin` role for testing offer CRUD.

## See what is stored

Open PostgreSQL's command-line client inside the database container:

```powershell
docker compose -f compose.yaml -f compose.local.yaml exec db psql -U eshche_est -d eshche_est
```

Useful SQL commands:

```sql
\dt

SELECT id, email, name, city, role, verified_at, created_at
FROM users
ORDER BY created_at DESC;

SELECT id, user_id, expires_at, revoked_at, created_at
FROM user_sessions
ORDER BY created_at DESC;

SELECT id, title, merchant, price_kopecks, quantity, status
FROM offers
ORDER BY created_at DESC;
```

Exit `psql` with `\q`.

You can also connect with DBeaver, DataGrip, or pgAdmin:

```text
Host: 127.0.0.1
Port: 5433
Database: eshche_est
User: eshche_est
Password: local-development-password
```

These credentials are deliberately local-only and must never be used on the
VPS.

## What writes information to PostgreSQL

The frontend never writes SQL directly. It sends JSON to the Go API, and the Go
repository layer runs parameterized database queries:

| Action | Endpoint | Tables changed |
|---|---|---|
| Request registration code | `POST /api/v1/auth/register/request` | `email_challenges` |
| Verify registration | `POST /api/v1/auth/register/verify` | `users`, `email_challenges`, `user_sessions` |
| Verify login | `POST /api/v1/auth/login/verify` | `email_challenges`, `user_sessions` |
| Edit profile | `PATCH /api/v1/users/me` | `users` |
| Logout | `POST /api/v1/auth/logout` | `user_sessions` |
| Create/update/delete an offer | `/api/v1/offers` | `offers` |

For a new kind of information, add it explicitly rather than creating an
unrestricted “store anything” endpoint:

1. Add a new versioned SQL migration under `apps/api/migrations`.
2. Add the Go model and parameterized queries under `apps/api/internal/store`.
3. Add validated HTTP handlers under `apps/api/internal/httpapi`.
4. Add the endpoint and schema to `apps/api/api/openapi.yaml`.
5. Call that endpoint from a Nuxt composable.

The API runs all pending migrations automatically when it starts.

## Stop, restart, and reset

Stop the containers without losing database data:

```powershell
docker compose -f compose.yaml -f compose.local.yaml down
```

Start them later with the same data:

```powershell
docker compose -f compose.yaml -f compose.local.yaml up -d
```

To permanently delete the local database and start from an empty migration,
use the following destructive command only when that is what you want:

```powershell
docker compose -f compose.yaml -f compose.local.yaml down --volumes
```

## Use the PostgreSQL database on your VPS

Local Compose and VPS Compose use the same migrations and API code. For the VPS,
do not include `compose.local.yaml`. Copy `.env.docker.example` to `.env`, set
the real `DATABASE_URL`, `OTP_PEPPER`, SMTP credentials, and admin email, then:

```bash
docker compose up -d --build
```

All API writes will then go to the database identified by `DATABASE_URL`.
