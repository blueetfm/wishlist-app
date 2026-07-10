"use client";

import { useState } from "react";

interface CommentFormProps {
  wishlistId?: string;
  shareToken?: string;
}

export function CommentForm({ wishlistId, shareToken }: CommentFormProps) {
  const [content, setContent] = useState("");
  const [authorName, setAuthorName] = useState("");
  const [status, setStatus] = useState<string | null>(null);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus("Sending comment...");
    setTimeout(() => setStatus("Comment submitted (placeholder)."), 500);
    setContent("");
  };

  return (
    <form onSubmit={handleSubmit} className="mt-6 space-y-4 rounded-3xl border border-slate-200 bg-slate-50 p-6">
      <div>
        <label className="block text-sm font-medium text-slate-700">Name</label>
        <input
          value={authorName}
          onChange={(event) => setAuthorName(event.target.value)}
          className="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
          placeholder="Optional"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-slate-700">Comment</label>
        <textarea
          value={content}
          onChange={(event) => setContent(event.target.value)}
          className="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-slate-900 outline-none focus:border-slate-400"
          placeholder="Write a comment..."
          rows={4}
        />
      </div>

      <button
        type="submit"
        className="inline-flex rounded-2xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white hover:bg-slate-800"
      >
        Post comment
      </button>

      {status ? <p className="text-sm text-slate-500">{status}</p> : null}
    </form>
  );
}
