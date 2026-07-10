"use client";

import { useState } from "react";

export function ScraperForm() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<string | null>(null);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setResult(`Extracted metadata for ${url} (placeholder)`);
  };

  return (
    <section className="rounded-3xl bg-white p-8 shadow-lg">
      <h2 className="text-2xl font-semibold">Product extractor</h2>
      <p className="mt-2 text-slate-600">Paste a product URL to fetch metadata from the scraper API.</p>
      <form onSubmit={handleSubmit} className="mt-6 grid gap-4 sm:grid-cols-[1fr_auto]">
        <input
          value={url}
          onChange={(event) => setUrl(event.target.value)}
          placeholder="https://example.com/product"
          className="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
        />
        <button className="rounded-2xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white hover:bg-slate-800">
          Extract
        </button>
      </form>
      {result ? <p className="mt-4 text-sm text-slate-500">{result}</p> : null}
    </section>
  );
}
