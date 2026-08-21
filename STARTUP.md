# ЕщёЕсть startup guide

This project supports three startup modes. Choose one mode at a time:

1. **Fully local** — Nuxt, Go API, and PostgreSQL all run in Docker locally.
2. **Local app + VPS PostgreSQL** — Nuxt and API run locally while PostgreSQL remains on your server.
3. **VPS deployment** — Nuxt and API run on the VPS and use PostgreSQL already running there.

Do not include `compose.local.yaml` when you want to use the server database:
that file intentionally overrides `DATABASE_URL` with the local PostgreSQL
container.

## Requirements

- Docker Desktop on Windows, or Docker Engine with Compose on Linux.
- Ports `3000` and `8080` available for the application.
- Port `5433` available only for the fully local PostgreSQL mode.

Check Docker:

```powershell
docker version
docker compose version
```

## Mode 1: fully local

This is the recommended mode for normal development. It does not need your VPS
or an SMTP account.

### Start

Run from the repository root:

```powershell
docker compose -f compose.yaml -f compose.local.yaml up -d --build
docker compose -f compose.yaml -f compose.local.yaml ps
```

Open:

- Application: `http://localhost:3000`
- Registration: `http://localhost:3000/register`
- Login: `http://localhost:3000/login`
- API health: `http://127.0.0.1:8080/healthz`

Local verification codes are generated randomly but are not emailed. The API
returns the code as `devCode`; the frontend displays and fills it automatically.

The local PostgreSQL database is available to database tools at:

```text
Host: 127.0.0.1
Port: 5433
Database: eshche_est
User: eshche_est
Password: local-development-password
```

These credentials are only for the local container.

### Inspect stored data

```powershell
docker compose -f compose.yaml -f compose.local.yaml exec db psql -U eshche_est -d eshche_est
```

Inside `psql`:

```sql
\dt

SELECT id, email, name, city, role, created_at
FROM users
ORDER BY created_at DESC;

SELECT id, user_id, expires_at, revoked_at
FROM user_sessions
ORDER BY created_at DESC;

SELECT id, title, merchant, price_kopecks, quantity, status
FROM offers
ORDER BY created_at DESC;
```

Exit with:

```sql
\q
```

### Stop and restart

Stop the services while preserving PostgreSQL data:

```powershell
docker compose -f compose.yaml -f compose.local.yaml stop
```

Start the same containers again:

```powershell
docker compose -f compose.yaml -f compose.local.yaml start
```

Remove the containers while preserving PostgreSQL data:

```powershell
docker compose -f compose.yaml -f compose.local.yaml down
```

Recreate them later with the same stored data:

```powershell
docker compose -f compose.yaml -f compose.local.yaml up -d
```

The data persists in the named volume
`eshche-est-local_postgres_local_data`.

To permanently remove the local database, use this destructive command only
when you intentionally want an empty database:

```powershell
docker compose -f compose.yaml -f compose.local.yaml down --volumes
```

## Mode 2: local app with PostgreSQL on your VPS

The safest approach is an SSH tunnel. PostgreSQL can remain bound to the VPS
loopback interface and port `5432` does not need to be exposed publicly.

### 1. Stop the fully local stack

This keeps its local database volume:

```powershell
docker compose -f compose.yaml -f compose.local.yaml down
```

### 2. Open an SSH tunnel

Keep this command running in a separate terminal:

```powershell
ssh -N -L 15432:127.0.0.1:5432 VPS_USER@VPS_HOST
```

The tunnel maps local port `15432` to PostgreSQL port `5432` on the VPS.

### 3. Create the application environment

```powershell
Copy-Item .env.docker.example .env
```

Edit `.env`:

```dotenv
WEB_PORT=3000
API_PORT=8080
APP_VERSION=latest

DATABASE_URL=postgres://DB_USER:URL_ENCODED_DB_PASSWORD@host.docker.internal:15432/DB_NAME?sslmode=disable
MIGRATE_ON_START=true

# Generate a private value of at least 32 characters.
OTP_PEPPER=REPLACE_WITH_A_LONG_RANDOM_SECRET
OTP_EXPIRY=10m
OTP_COOLDOWN=1m
OTP_MAX_ATTEMPTS=5
SESSION_TTL=720h
SESSION_COOKIE_SECURE=false

# No SMTP is required for a private local test.
EMAIL_DELIVERY=development
EXPOSE_DEV_CODES=true

# An account registered with this address receives admin offer permissions.
ADMIN_EMAILS=YOUR_EMAIL@example.ru
```

`DB_USER`, `DB_PASSWORD`, and `DB_NAME` must be real PostgreSQL credentials from
your VPS. URL-encode reserved password characters such as `@`, `:`, `/`, `?`,
`#`, and `%`.

The `.env` file contains secrets and must not be committed. It is already listed
in `.gitignore`.

### 4. Start only Nuxt and the API

Do not add `compose.local.yaml` here:

```powershell
docker compose up -d --build
docker compose ps
docker compose logs --tail=100 api web
```

Open `http://localhost:3000`. Registration, login, profiles, sessions, and API
offer changes will now be written to the VPS database identified by
`DATABASE_URL`.

### 5. Stop

```powershell
docker compose down
```

Stop the SSH tunnel with `Ctrl+C` in its terminal.

## Mode 3: application and PostgreSQL on the VPS

Use this mode when other people need to open the application from their devices.

### 1. Copy and configure the environment

On the VPS, from the repository root:

```bash
cp .env.docker.example .env
openssl rand -hex 32
```

Edit `.env` and place the generated value in `OTP_PEPPER`.

If PostgreSQL runs directly on the same VPS host:

```dotenv
DATABASE_URL=postgres://DB_USER:URL_ENCODED_DB_PASSWORD@host.docker.internal:5432/DB_NAME?sslmode=disable
```

The PostgreSQL server must listen on an interface reachable from Docker's bridge
network, and `pg_hba.conf` must allow the specific Docker subnet. Keep port
`5432` blocked from the public internet.

If PostgreSQL is on a different server or provided as a managed database, use
its private hostname and TLS:

```dotenv
DATABASE_URL=postgres://DB_USER:URL_ENCODED_DB_PASSWORD@PRIVATE_DB_HOST:5432/DB_NAME?sslmode=require
```

### 2. Choose email mode

For a private preview without SMTP:

```dotenv
EMAIL_DELIVERY=development
EXPOSE_DEV_CODES=true
SESSION_COOKIE_SECURE=false
```

Anyone who can reach that preview can see their returned verification code, so
do not treat development mode as proof of email ownership.

For real email verification:

```dotenv
EMAIL_DELIVERY=smtp
EXPOSE_DEV_CODES=false
SMTP_HOST=YOUR_SMTP_HOST
SMTP_PORT=587
SMTP_USER=YOUR_SMTP_USER
SMTP_PASSWORD=YOUR_SMTP_OR_APP_PASSWORD
SMTP_FROM=ЕщёЕсть <YOUR_AUTHORIZED_SENDER@example.ru>
SESSION_COOKIE_SECURE=true
```

SMTP mode requires STARTTLS on port `587`. `SMTP_FROM` normally needs to be an
address authorized by the SMTP provider.

### 3. Build and start

```bash
docker compose up -d --build
docker compose ps
docker compose logs --tail=100 api web
```

Check on the VPS:

```bash
curl --fail http://127.0.0.1:8080/healthz
curl --fail http://127.0.0.1:3000/
curl --fail http://127.0.0.1:3000/api/v1/offers
```

The API port is bound to VPS loopback. Only expose the Nuxt web port, or put
Caddy/Nginx in front of `127.0.0.1:3000` and expose HTTPS ports `80/443`.

### 4. Stop or update

Stop:

```bash
docker compose down
```

Update after pulling new code:

```bash
docker compose build
docker compose up -d
docker compose logs --tail=100 api web
```

## How database storage works

The browser never connects directly to PostgreSQL:

```text
Browser → Nuxt /api proxy → Go API → PostgreSQL
```

The API uses parameterized queries and automatically runs pending migrations
from `apps/api/migrations` during startup.

Current stored information:

| User action | PostgreSQL tables |
|---|---|
| Request registration/login code | `email_challenges` |
| Complete registration | `users`, `email_challenges`, `user_sessions` |
| Complete login | `email_challenges`, `user_sessions` |
| Update profile | `users` |
| Logout | `user_sessions` |
| Create/update/delete offer through the API | `offers` |

The frontend offer catalogue still uses mock data. Offer records created through
the API are stored in PostgreSQL, but connecting the browse/discover UI to those
records is a separate frontend integration step.

## Troubleshooting

### See service status

```powershell
docker compose -f compose.yaml -f compose.local.yaml ps
```

For server-database mode, omit `-f compose.local.yaml`:

```powershell
docker compose ps
```

### Read logs

```powershell
docker compose -f compose.yaml -f compose.local.yaml logs --tail=200 api db web
```

### API reports that PostgreSQL is unavailable

Check:

- the `DATABASE_URL` username, password, host, port, and database name;
- whether the SSH tunnel is still running;
- whether PostgreSQL accepts the application user;
- PostgreSQL `listen_addresses` and `pg_hba.conf` when containers connect to the VPS host;
- firewall rules;
- `sslmode=require` versus `sslmode=disable`.

### Login says the code is invalid

- Request a new code; codes expire after ten minutes by default.
- Use the latest code for that email and purpose.
- In local development, use the displayed `devCode`.
- In SMTP mode, use the code from the newest email.

### Rebuild after source changes

Local mode:

```powershell
docker compose -f compose.yaml -f compose.local.yaml up -d --build
```

Server-database mode:

```powershell
docker compose up -d --build
```
