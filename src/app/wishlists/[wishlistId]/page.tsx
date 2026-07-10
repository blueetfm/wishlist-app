import { CommentForm } from "@/app/components/comments/comment-form";
import { CommentThread } from "@/app/components/comments/comment-thread";
import { ScraperForm } from "@/app/components/scraper/url-extract-form";
import { WishlistItemList } from "@/app/components/wishlist/wishlist-item-list";

interface WishlistDetailPageProps {
  params: {
    wishlistId: string;
  };
}

export default function WishlistDetailPage({ params: { wishlistId } }: WishlistDetailPageProps) {
  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10">
      <div className="mx-auto max-w-6xl space-y-8">
        <section className="rounded-3xl bg-white p-8 shadow-lg">
          <h1 className="text-3xl font-semibold">Wishlist details</h1>
          <p className="mt-2 text-slate-600">Viewing wishlist <span className="font-medium">{wishlistId}</span>.</p>
        </section>

        <WishlistItemList />

        <ScraperForm />

        <section className="rounded-3xl bg-white p-8 shadow-lg">
          <h2 className="text-2xl font-semibold">Comments</h2>
          <CommentThread />
          <CommentForm wishlistId={wishlistId} />
        </section>
      </div>
    </main>
  );
}
