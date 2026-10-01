# Wishlist Frontend

Next.js (App Router) frontend for the Wishlist app. Supabase Auth (Google
sign-in) and app shell only — wishlist UI is built against the
[backend API](../backend/docs/swagger.yaml).

## Stack

Next.js 15 · React 18 · TypeScript · Tailwind CSS · Radix UI · Supabase

## Setup

```bash
npm install
npm run dev
```

Create `frontend/.env.local`:

```bash
NEXT_PUBLIC_SUPABASE_URL=https://<your-project>.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=<your-anon-key>
```

Open [http://localhost:3000](http://localhost:3000).
