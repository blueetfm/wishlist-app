import React from "react";
import { ImageOff } from "lucide-react";

interface ProductImageProps {
  /** User-uploaded photo (Item.image_url). Takes precedence when present. */
  imageUrl?: string | null;
  /** Open Graph image scraped from the product link (Item.embed_data.image_url). */
  embedImageUrl?: string | null;
  alt: string;
  className?: string;
}

/** Resolves a product's photo for use as GiftBox's children: upload > embed link > placeholder. */
export function ProductImage({
  imageUrl,
  embedImageUrl,
  alt,
  className = "",
}: ProductImageProps) {
  const src = imageUrl || embedImageUrl || null;

  if (!src) {
    return (
      <div
        className={`flex items-center justify-center w-full h-full min-h-[160px] text-muted-foreground ${className}`}
      >
        <ImageOff className="w-10 h-10" strokeWidth={1.5} />
      </div>
    );
  }

  return (
    // Product photos come from arbitrary user uploads and embed domains, so
    // next/image's fixed domain allowlist isn't a good fit here.
    // eslint-disable-next-line @next/next/no-img-element
    <img
      src={src}
      alt={alt}
      className={`w-full h-full object-cover rounded-xl ${className}`}
    />
  );
}

export default ProductImage;
