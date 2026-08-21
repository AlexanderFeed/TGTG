# Deploy ЕщёЕсть web and API to a VPS with Docker

The Compose stack runs the Nuxt frontend and Go API. The browser reaches only
Nuxt; `/api/*` is proxied over the private Compose network to the API. Port 8080
is bound to VPS loopback for diagnostics and is not publicly exposed.

The API uses your existing PostgreSQL instance and runs versioned migrations at
startup. It does not create or delete the database itself.

## VPS requirements

- Linux VPS with Docker Engine and Docker Compose.
- PostgreSQL 14+ with an empty database and a dedicated login for this app.
- SMTP submission endpoint on port 587 with STARTTLS.
- Project copied or cloned onto the VPS.
- TCP port 3000 open for an IP preview, or ports 80/443 through Caddy/Nginx.

## 1. Prepare PostgreSQL

Create a database and least-privilege application login. The login currently
needs schema migration rights because the API runs migrations on startup.

If PostgreSQL runs directly on the VPS, it must listen on an address reachable
from Docker's bridge network, and `pg_hba.conf` must permit the Docker subnet.
Keep port 5432 blocked from the public internet. In `.env`, use
`host.docker.internal`, which Compose maps to the VPS host:

```dotenv
DATABASE_URL=postgres://eshche_est:URL_ENCODED_PASSWORD@host.docker.internal:5432/eshche_est?sslmode=disable
```

Use `sslmode=require` when PostgreSQL is remote or TLS is configured. If the
database is another Docker service, attach it to the same network and use its
Compose service name instead.

## 2. Configure secrets and email

From the repository root:

```bash
cp .env.docker.example .env
openssl rand -hex 32
```

Put the generated value in `OTP_PEPPER`, then set `DATABASE_URL`, the SMTP
values, and `ADMIN_EMAILS`. An email listed in `ADMIN_EMAILS` becomes an admin
when it registers and can use offer write endpoints.

For a plain `http://IP:3000` preview, keep:

```dotenv
SESSION_COOKIE_SECURE=false
```

Behind an HTTPS domain, set it to `true`. Production email verification uses:

```dotenv
EMAIL_DELIVERY=smtp
EXPOSE_DEV_CODES=false
SMTP_HOST=smtp.example.ru
SMTP_PORT=587
SMTP_USER=mailer@example.ru
SMTP_PASSWORD=your-secret
SMTP_FROM=ЕщёЕсть <mailer@example.ru>
```

`EMAIL_DELIVERY=development` with `EXPOSE_DEV_CODES=true` returns the code in
the API response. Use that combination only on a private local environment—it
does not prove email ownership.

## 3. Build and start

```bash
docker compose build
docker compose up -d
docker compose ps
docker compose logs --tail=100 api web
```

Check both services locally on the VPS:

```bash
curl --fail http://127.0.0.1:8080/healthz
curl --fail http://127.0.0.1:3000/
curl --fail http://127.0.0.1:3000/api/v1/offers
```

Then open:

```text
http://YOUR_VPS_PUBLIC_IP:3000/
```

If local checks succeed but another device cannot connect, allow only the web
port in both firewalls. For Ubuntu UFW:

```bash
sudo ufw allow 3000/tcp
sudo ufw status
```

Do not open API port 8080 or PostgreSQL port 5432 publicly.

## HTTPS domain

Put Caddy or Nginx in front of `127.0.0.1:3000`, issue a TLS certificate, set
`SESSION_COOKIE_SECURE=true`, and expose only ports 80/443. The Go API remains
behind Nuxt, so no separate API domain or CORS configuration is needed.

## Update and inspect

```bash
docker compose build
docker compose up -d
docker compose ps
docker compose logs -f api web
```

Useful checks:

```bash
docker compose restart api web
docker inspect --format='{{json .State.Health}}' eshche-est-api-1
docker inspect --format='{{json .State.Health}}' eshche-est-web-1
```

## Backup and rollback

Back up PostgreSQL before deploying a migration. The first migration is
additive, but future rollbacks must account for both code and schema versions.
For a code rollback, check out the known-good release and rebuild:

```bash
git checkout YOUR_LAST_GOOD_COMMIT
docker compose build
docker compose up -d
```

Use `APP_VERSION` with a commit SHA or release number for identifiable images.
