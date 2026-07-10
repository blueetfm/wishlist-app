import { CommentForm } from "@/app/components/comments/comment-form";
import { WishlistItemList } from "@/app/components/wishlist/wishlist-item-list";

interface SharedWishlistPageProps {
  params: {
    shareToken: string;
  };
}

export default function SharedWishlistPage({ params: { shareToken } }: SharedWishlistPageProps) {
  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10">
      <div className="mx-auto max-w-6xl space-y-8">
        <section className="rounded-3xl bg-white p-8 shadow-lg">
          <h1 className="text-3xl font-semibold">Shared wishlist</h1>
          <p className="mt-2 text-slate-600">Sharing with token <span className="font-medium">{shareToken}</span>.</p>
        </section>

        <WishlistItemList guestView />

        <section className="rounded-3xl bg-white p-8 shadow-lg">
          <h2 className="text-2xl font-semibold">Leave a comment</h2>
          <p className="mt-2 text-slate-600">Guests can comment anonymously or leave a name.</p>
          <CommentForm shareToken={shareToken} />
        </section>
      </div>
    </main>
  );
}
