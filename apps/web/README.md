# ЕщёЕсть customer web

Nuxt 4 / Vue 3 frontend prototype for a Russian food-rescue marketplace. It currently uses local mock data and browser cookies; no backend or payment provider is connected.

## Run locally

Requirements: Node.js 20+ and npm.

```powershell
npm install
npm run dev
```

Open `http://localhost:3000`.

## Demo flows

- Guest home: `/`
- Login: `/login` — enter a valid-looking Russian phone number, then any four digits
- Registration: `/register`
- Nearby map/list: `/discover`
- Browse and filters: `/browse`
- Offer detail and reservation prototype: `/offers/bakery-evening`
- Pickup order and future delivery concept: `/delivery`
- Profile and favorites: `/profile`

Authentication, favorites, and profile changes are stored only in cookies. The delivery tab is intentionally marked as a future concept because the MVP in `PLAN.md` remains pickup-only.

## Verification

```powershell
npm run typecheck
npm run build
```

The current implementation passes both commands.
