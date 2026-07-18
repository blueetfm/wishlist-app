"use client";

import { useState } from "react";
import { buildUrl } from "@/lib/utils";

const DEFAULT_API_BASE_URL = process.env.NEXT_PUBLIC_WISHLIST_API_BASE_URL ?? "http://localhost:4000";

async function readResponseBody(response: Response) {
  const text = await response.text();

  if (!text) {
    return "(empty response)";
  }

  try {
    return JSON.stringify(JSON.parse(text), null, 2);
  } catch {
    return text;
  }
}

export function WishlistApiTester() {
  const [baseUrl, setBaseUrl] = useState(DEFAULT_API_BASE_URL);
  const [wishlistId, setWishlistId] = useState("abc123");
  const [shareToken, setShareToken] = useState("example-token");
  const [itemId, setItemId] = useState("item-1");
  const [title, setTitle] = useState("Weekend wishlist");
  const [description, setDescription] = useState("Things to test in the API");
  const [itemName, setItemName] = useState("Noise-cancelling headphones");
  const [productUrl, setProductUrl] = useState("https://example.com/product");
  const [guestName, setGuestName] = useState("Avery");
  const [message, setMessage] = useState("Please reserve this item.");
  const [loadingLabel, setLoadingLabel] = useState<string | null>(null);
  const [result, setResult] = useState<{
    label: string;
    method: string;
    url: string;
    status: number | null;
    body: string;
  } | null>(null);

  const runRequest = async (label: string, path: string, init: RequestInit = {}) => {
    const url = buildUrl(baseUrl, path);

    setLoadingLabel(label);
    setResult(null);

    try {
      const response = await fetch(url, {
        headers: {
          "Content-Type": "application/json",
          ...(init.headers || {}),
        },
        ...init,
      });

      const body = await readResponseBody(response);

      setResult({
        label,
        method: init.method || "GET",
        url,
        status: response.status,
        body,
      });
    } catch (error) {
      setResult({
        label,
        method: init.method || "GET",
        url,
        status: null,
        body: error instanceof Error ? error.message : "Unknown request error",
      });
    } finally {
      setLoadingLabel(null);
    }
  };

  return (
    <section className="rounded-3xl bg-white/70 p-8 shadow-lg backdrop-blur-sm">
      <div className="flex flex-col gap-3 border-b border-slate-200 pb-6 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <h2 className="text-2xl font-semibold text-slate-900">Wishlist API tester</h2>
          <p className="mt-2 max-w-2xl text-sm text-slate-600">
            Point this at the Fastify backend by default, or override it if you want to hit a different backend URL.
          </p>
        </div>
        <label className="space-y-2 text-sm font-medium text-slate-700 lg:min-w-[22rem]">
          API base URL
          <input
            value={baseUrl}
            onChange={(event) => setBaseUrl(event.target.value)}
            className="w-full rounded-2xl border border-slate-200 bg-white/70 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
            placeholder="/api"
          />
        </label>
      </div>

      <div className="mt-6 grid gap-4 lg:grid-cols-2">
        <article className="rounded-3xl border border-slate-200 bg-white/55 p-5">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-400">Read</p>
              <h3 className="mt-1 text-lg font-semibold text-slate-900">List wishlists</h3>
            </div>
            <button
              type="button"
              onClick={() => void runRequest("List wishlists", "/api/wishlists")}
              className="rounded-xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
              disabled={loadingLabel === "List wishlists"}
            >
              {loadingLabel === "List wishlists" ? "Loading..." : "GET"}
            </button>
          </div>
          <p className="mt-3 text-sm text-slate-600">Fetch the public wishlist collection.</p>
        </article>

        <article className="rounded-3xl border border-slate-200 bg-white/55 p-5">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-400">Auth</p>
              <h3 className="mt-1 text-lg font-semibold text-slate-900">My wishlists</h3>
            </div>
            <button
              type="button"
              onClick={() => void runRequest("My wishlists", "/api/wishlists/me")}
              className="rounded-xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
              disabled={loadingLabel === "My wishlists"}
            >
              {loadingLabel === "My wishlists" ? "Loading..." : "GET"}
            </button>
          </div>
          <p className="mt-3 text-sm text-slate-600">Exercise the authenticated wishlist endpoint.</p>
        </article>

        <article className="rounded-3xl border border-slate-200 bg-white/55 p-5">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-400">Read</p>
              <h3 className="mt-1 text-lg font-semibold text-slate-900">Load wishlist items</h3>
            </div>
            <button
              type="button"
              onClick={() => void runRequest("Load wishlist items", `/api/wishlists/${wishlistId}/items`)}
              className="rounded-xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
              disabled={loadingLabel === "Load wishlist items"}
            >
              {loadingLabel === "Load wishlist items" ? "Loading..." : "GET"}
            </button>
          </div>
          <label className="mt-4 block space-y-2 text-sm font-medium text-slate-700">
            Wishlist ID
            <input
              value={wishlistId}
              onChange={(event) => setWishlistId(event.target.value)}
              className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
            />
          </label>
        </article>

        <article className="rounded-3xl border border-slate-200 bg-white/55 p-5">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-400">Write</p>
              <h3 className="mt-1 text-lg font-semibold text-slate-900">Create wishlist</h3>
            </div>
            <button
              type="button"
              onClick={() =>
                void runRequest("Create wishlist", "/api/wishlists/new", {
                  method: "POST",
                  body: JSON.stringify({
                    title,
                    description,
                  }),
                })
              }
              className="rounded-xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
              disabled={loadingLabel === "Create wishlist"}
            >
              {loadingLabel === "Create wishlist" ? "Sending..." : "POST"}
            </button>
          </div>
          <div className="mt-4 grid gap-3">
            <label className="space-y-2 text-sm font-medium text-slate-700">
              Title
              <input
                value={title}
                onChange={(event) => setTitle(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
            <label className="space-y-2 text-sm font-medium text-slate-700">
              Description
              <input
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
          </div>
        </article>

        <article className="rounded-3xl border border-slate-200 bg-white/55 p-5 lg:col-span-2">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-400">Write</p>
              <h3 className="mt-1 text-lg font-semibold text-slate-900">Add wishlist item</h3>
            </div>
            <button
              type="button"
              onClick={() =>
                void runRequest("Add wishlist item", `/wishlists/${wishlistId}/items`, {
                  method: "POST",
                  body: JSON.stringify({
                    name: itemName,
                    productUrl,
                  }),
                })
              }
              className="rounded-xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
              disabled={loadingLabel === "Add wishlist item"}
            >
              {loadingLabel === "Add wishlist item" ? "Sending..." : "POST"}
            </button>
          </div>
          <div className="mt-4 grid gap-3 md:grid-cols-3">
            <label className="space-y-2 text-sm font-medium text-slate-700">
              Item name
              <input
                value={itemName}
                onChange={(event) => setItemName(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
            <label className="space-y-2 text-sm font-medium text-slate-700 md:col-span-2">
              Product URL
              <input
                value={productUrl}
                onChange={(event) => setProductUrl(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
          </div>
        </article>

        <article className="rounded-3xl border border-slate-200 bg-white/55 p-5 lg:col-span-2">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.2em] text-slate-400">Public</p>
              <h3 className="mt-1 text-lg font-semibold text-slate-900">Claim shared item</h3>
            </div>
            <button
              type="button"
              onClick={() =>
                void runRequest(
                  "Claim shared item",
                  `/api/wishlists/shared/${shareToken}/items/${itemId}/claim`,
                  {
                    method: "POST",
                    body: JSON.stringify({
                      guestName,
                      message,
                    }),
                  },
                )
              }
              className="rounded-xl bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-800"
              disabled={loadingLabel === "Claim shared item"}
            >
              {loadingLabel === "Claim shared item" ? "Sending..." : "POST"}
            </button>
          </div>
          <div className="mt-4 grid gap-3 md:grid-cols-3">
            <label className="space-y-2 text-sm font-medium text-slate-700">
              Share token
              <input
                value={shareToken}
                onChange={(event) => setShareToken(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
            <label className="space-y-2 text-sm font-medium text-slate-700">
              Item ID
              <input
                value={itemId}
                onChange={(event) => setItemId(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
            <label className="space-y-2 text-sm font-medium text-slate-700">
              Guest name
              <input
                value={guestName}
                onChange={(event) => setGuestName(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
            <label className="space-y-2 text-sm font-medium text-slate-700 md:col-span-3">
              Message
              <input
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                className="w-full rounded-2xl border border-slate-200 bg-white/75 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
              />
            </label>
          </div>
        </article>
      </div>

      <div className="mt-6 rounded-3xl border border-slate-200 bg-slate-950/85 p-6 text-slate-100 backdrop-blur-sm">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <p className="text-sm uppercase tracking-[0.2em] text-slate-400">Latest response</p>
            <p className="mt-1 text-lg font-semibold">{result?.label ?? "Nothing sent yet"}</p>
          </div>
          {result ? (
            <div className="text-right text-sm text-slate-300">
              <p>
                {result.method} {result.url}
              </p>
              <p>{result.status === null ? "Request failed" : `HTTP ${result.status}`}</p>
            </div>
          ) : null}
        </div>

        <pre className="mt-4 overflow-x-auto rounded-2xl bg-slate-900 p-4 text-sm leading-6 text-slate-100">
{result ? result.body : "Run a request to see the response body here."}
        </pre>
      </div>
    </section>
  );
}