import Link from "next/link";

export default function Page() {
  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10">
      <div className="mx-auto max-w-5xl rounded-3xl bg-white p-10 shadow-lg">
        <h1 className="text-4xl font-bold tracking-tight">Wishlist App</h1>
        <p className="mt-4 text-lg text-slate-600">
          Manage your personal wishlists, share them with friends, and track claims and comments.
        </p>

        <div className="mt-10 grid gap-4 sm:grid-cols-2">
          <Link
            href="/wishlists"
            className="rounded-2xl border border-slate-200 bg-slate-900 px-6 py-5 text-white transition hover:bg-slate-800"
          >
            My Wishlists
          </Link>
          <Link
            href="/swagger"
            className="rounded-2xl border border-slate-200 bg-white px-6 py-5 text-slate-900 transition hover:bg-slate-50"
          >
            API / Swagger
          </Link>
        </div>
      </div>
    </main>
  );
}
