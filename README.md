# Wishlist App

A monorepo for the Wishlist app: a Next.js frontend and a Go backend sharing
a Supabase project. Create shareable wishlists, claim items 1-to-1, pool
interest to split costs, embed Open Graph link previews, and comment on
wishlists/items.

```
frontend/    Next.js frontend — Supabase Auth, UI
backend/     Go REST API — wishlists, items, claims, comments, embeds
```

## Getting started

```bash
cd frontend && npm install && npm run dev   # http://localhost:3000
cd backend && go mod tidy && go run ./cmd/api   # http://localhost:8080
```

See [frontend/README.md](frontend/README.md) and
[backend/README.md](backend/README.md) for required `.env` variables.

## API

Full REST contract: [backend/docs/swagger.yaml](backend/docs/swagger.yaml)
(OpenAPI 3.0).
