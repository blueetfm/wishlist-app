interface CommentThreadProps {
  comments?: Array<{
    id: string;
    authorName?: string | null;
    content: string;
  }>;
}

export function CommentThread({
  comments = [
    { id: "c1", authorName: "Avery", content: "Really love this wishlist!" },
    { id: "c2", authorName: null, content: "Anonymous guest claim request." },
  ],
}: CommentThreadProps) {
  return (
    <div className="space-y-4 py-4">
      {comments.map((comment) => (
        <article key={comment.id} className="rounded-3xl border border-slate-200 bg-slate-50 p-5">
          <div className="flex items-center justify-between gap-4">
            <span className="text-sm font-semibold text-slate-900">{comment.authorName || "Guest"}</span>
            <span className="text-xs uppercase tracking-[0.2em] text-slate-400">Comment</span>
          </div>
          <p className="mt-3 text-slate-700">{comment.content}</p>
        </article>
      ))}
    </div>
  );
}
