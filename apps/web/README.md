# ЕщёЕсть customer web

Nuxt 4 / Vue 3 customer frontend for a Russian food-rescue marketplace. Email authentication and profile edits use the Go/PostgreSQL backend through a same-origin Nuxt proxy. Offers, favorites, reservations, orders, and payments still use prototype data/state.

Guests can use the main navigation, browse offers, open offer details, and grant
browser geolocation permission. Authentication is requested when they try to
reserve an offer or open private profile/order data. Geolocation works on
`localhost` or HTTPS; the current illustrative map does not yet calculate real
distances from the captured coordinates.

## Run locally

Requirements: Node.js 22+, npm, and the API from `apps/api` on port 8080.

```powershell
npm install
npm run dev
```

Open `http://localhost:3000`.

## Demo flows

- Guest home: `/`
- Login: `/login` — request and verify a six-digit email code
- Registration: `/register` — verify email, create the PostgreSQL user, and start a session
- Nearby map/list: `/discover`
- Browse and filters: `/`
- Offer detail and reservation prototype: `/offers/bakery-evening`
- Pickup order and future delivery concept: `/delivery`
- Profile and favorites: `/profile`

The opaque session token is stored in an HttpOnly cookie; session and user records live in PostgreSQL. Favorites remain in a browser cookie. The delivery tab is intentionally marked as a future concept because the MVP in `PLAN.md` remains pickup-only.

## Verification

```powershell
npm run typecheck
npm run build
```

The current implementation passes both commands.

## Docker

For a complete local PostgreSQL/API/Nuxt environment, run from the repository
root:

```powershell
docker compose -f compose.yaml -f compose.local.yaml up -d --build
```

The frontend is published at `http://localhost:3000`. See
`LOCAL_DEVELOPMENT.md` for database inspection and persistence commands, and
`DEPLOY.md` for VPS deployment, firewall, update, and rollback instructions.
