interface WishlistItemProps {
  name: string;
  price?: number | null;
  productUrl?: string | null;
  imageUrl?: string | null;
}

export function WishlistItem({ name, price, productUrl, imageUrl }: WishlistItemProps) {
  return (
    <li className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex items-start gap-4">
        <div className="flex-1 space-y-2">
          <h3 className="text-lg font-semibold">{name}</h3>
          <p className="text-sm text-slate-500">{productUrl || "No product URL provided."}</p>
          {price != null ? <p className="text-sm font-medium text-slate-700">Price: ${price.toFixed(2)}</p> : null}
        </div>
        {imageUrl ? (
          <img src={imageUrl} alt={name} className="h-16 w-16 rounded-3xl object-cover" />
        ) : (
          <div className="flex h-16 w-16 items-center justify-center rounded-3xl bg-slate-100 text-slate-400">No image</div>
        )}
      </div>
    </li>
  );
}
