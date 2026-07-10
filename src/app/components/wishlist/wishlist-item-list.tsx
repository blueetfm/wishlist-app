import { WishlistItem } from "@/app/components/wishlist/wishlist-item";

interface WishlistItemListProps {
  guestView?: boolean;
}

const sampleItems = [
  {
    id: "item-1",
    name: "Noise canceling headphones",
    price: 199.99,
    productUrl: "https://example.com/headphones",
    imageUrl: null,
  },
  {
    id: "item-2",
    name: "Standing desk lamp",
    price: 49.95,
    productUrl: "https://example.com/lamp",
    imageUrl: null,
  },
];

export function WishlistItemList({ guestView = false }: WishlistItemListProps) {
  return (
    <section className="rounded-3xl bg-white p-8 shadow-lg">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-semibold">{guestView ? "Shared items" : "Wishlist items"}</h2>
          <p className="mt-2 text-slate-600">
            {guestView
              ? "Review the current shared wishlist and claim an item."
              : "Add or manage the items in this wishlist."}
          </p>
        </div>
        {guestView ? null : (
          <button className="rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white hover:bg-slate-800">
            Add item
          </button>
        )}
      </div>

      <ul className="mt-6 space-y-4">
        {sampleItems.map((item) => (
          <WishlistItem key={item.id} {...item} />
        ))}
      </ul>
    </section>
  );
}
