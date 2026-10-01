# Wishlist API (Backend)

Go REST API for the Wishlist app: wishlists, items, claims, split-interest,
Open Graph embeds, comments, and public shared wishlists.

## Stack

Go 1.26+ · `net/http` · Supabase (Postgres + Auth)

## Setup

```bash
go mod tidy
go run ./cmd/api
```

Create `backend/.env`:

```bash
PORT=8080
SUPABASE_URL=https://<your-project>.supabase.co
DATABASE_URL=postgres://<user>:<password>@<host>:5432/<database>
```

`SUPABASE_JWKS_URL` is optional (derived from `SUPABASE_URL` if unset).

## API

Full contract in [`docs/swagger.yaml`](docs/swagger.yaml). All endpoints
require a Supabase bearer token except `/shared/{shareToken}`.

## Build & test

```bash
go build ./cmd/api
go test ./...
```

