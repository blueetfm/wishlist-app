import Link from "next/link";
import { NewWishlistForm } from "@/app/components/wishlist/new-wishlist-form";

export default function CreateWishlistPage() {
  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10">
      <div className="mx-auto max-w-4xl space-y-8">
        <section className="rounded-3xl bg-white p-8 shadow-lg">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.25em] text-slate-400">Create wishlist</p>
              <h1 className="mt-2 text-3xl font-semibold text-slate-950">Start a new wishlist</h1>
              <p className="mt-2 max-w-2xl text-slate-600">
                This page POSTs directly to the wishlist API so you can test the create flow from the browser.
              </p>
            </div>
            <Link
              href="/wishlists"
              className="inline-flex rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm font-semibold text-slate-900 transition hover:bg-slate-50"
            >
              Back to wishlist API tester
            </Link>
          </div>
        </section>

        <NewWishlistForm />
      </div>
    </main>
  );
}