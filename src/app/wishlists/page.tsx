import Link from "next/link";
import { WishlistApiTester } from "@/app/components/wishlist/wishlist-api-tester";

export default function WishlistsPage() {
  return (
    <main className="min-h-screen bg-transparent px-6 py-10">
      <div className="mx-auto max-w-6xl space-y-8">
        <section className="rounded-3xl bg-white/70 p-8 shadow-lg backdrop-blur-sm">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.25em] text-slate-400">Testing ground</p>
              <h1 className="mt-2 text-3xl font-semibold text-slate-950">Wishlist API playground</h1>
              <p className="mt-2 max-w-3xl text-slate-600">
                Use this page to hit the wishlist endpoints directly from the browser and inspect the raw
                responses from your backend.
              </p>
            </div>
            <Link
              href="/"
              className="inline-flex rounded-2xl border border-slate-200 bg-white/80 px-4 py-3 text-sm font-semibold text-slate-900 transition hover:bg-white"
            >
              Back to home
            </Link>
          </div>
        </section>

        <WishlistApiTester />
      </div>
    </main>
  );
}
