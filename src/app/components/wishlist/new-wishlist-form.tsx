"use client";

import { useState } from "react";
import { useAuth } from "@/app/context/auth-context";
import { buildUrl } from "@/lib/utils";

const BACKEND_BASE_URL = process.env.NEXT_PUBLIC_WISHLIST_API_BASE_URL ?? "http://localhost:4000";

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

export function NewWishlistForm() {
  const { session } = useAuth();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState<string | null>(null);
  const [responseBody, setResponseBody] = useState<string | null>(null);
  const [createdWishlistId, setCreatedWishlistId] = useState<string | null>(null);
  const [createdShareToken, setCreatedShareToken] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const isTitleValid = title.trim().length > 0;

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!event.currentTarget.checkValidity()) {
      event.currentTarget.reportValidity();
      return;
    }

    const normalizedTitle = title.trim();

    setSubmitting(true);
    setStatus("Creating wishlist...");
    setResponseBody(null);
    setCreatedWishlistId(null);
    setCreatedShareToken(null);

    try {
      const accessToken = session?.access_token;

      if (!accessToken) {
        setStatus("Please sign in first. Missing Supabase access token.");
        return;
      }

      const response = await fetch(buildUrl(BACKEND_BASE_URL, "/api/wishlists/new"), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          title: normalizedTitle,
          description,
        }),
      });

      const bodyText = await readResponseBody(response);
      setResponseBody(bodyText);

      if (!response.ok) {
        setStatus(`Request failed with HTTP ${response.status}: ${bodyText}`);
        return;
      }

      const data = JSON.parse(bodyText) as {
        id?: string;
        userId?: string;
        title?: string;
        description?: string;
        createdAt?: string;
        shareToken?: string;
      };

      setCreatedWishlistId(data.id ?? null);
      setCreatedShareToken(data.shareToken ?? null);
      setStatus("Wishlist created successfully.");
      setTitle("");
      setDescription("");
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Failed to create wishlist.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="rounded-3xl bg-white p-8 shadow-lg">
      <h2 className="text-2xl font-semibold">Create a new wishlist</h2>
      <form onSubmit={handleSubmit} className="mt-6 space-y-4">
        <div>
          <label className="block text-sm font-medium text-slate-700">Title</label>
          <input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            required
            aria-invalid={!isTitleValid}
            className="mt-2 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
            placeholder="My birthday list"
          />
          {!isTitleValid ? <p className="mt-2 text-sm text-rose-600">Title is required.</p> : null}
        </div>

        <div>
          <label className="block text-sm font-medium text-slate-700">Description</label>
          <textarea
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            className="mt-2 w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
            placeholder="Ideas for the season"
            rows={4}
          />
        </div>

        <button
          type="submit"
          className="inline-flex rounded-2xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
          disabled={submitting || !isTitleValid}
        >
          {submitting ? "Creating..." : "Create wishlist"}
        </button>

        {status ? <p className="text-sm text-slate-500">{status}</p> : null}

        {createdWishlistId ? (
          <p className="text-sm text-emerald-700">
            Created wishlist ID: <span className="font-semibold">{createdWishlistId}</span>
          </p>
        ) : null}

        {createdShareToken ? (
          <p className="text-sm text-slate-500">
            Share token: <span className="font-semibold text-slate-900">{createdShareToken}</span>
          </p>
        ) : null}

        {responseBody ? (
          <pre className="overflow-x-auto rounded-2xl bg-slate-950 p-4 text-sm leading-6 text-slate-100">
            {responseBody}
          </pre>
        ) : null}
      </form>
    </section>
  );
}
