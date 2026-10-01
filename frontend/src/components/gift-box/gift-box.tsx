import React from "react";
import { GiftRibbon } from "./gift-ribbon";
import { ProductImage } from "./product-image";

interface GiftBoxProps {
  /** Overrides the built-in ProductImage rendering with arbitrary content. */
  children?: React.ReactNode;
  /** User-uploaded photo (Item.image_url) - takes precedence over embedImageUrl. */
  imageUrl?: string | null;
  /** Open Graph image scraped from the product link (Item.embed_data.image_url). */
  embedImageUrl?: string | null;
  imageAlt?: string;
  className?: string;
  boxClassName?: string;
  ribbonStrokeWidth?: number;
}

export default function GiftBox({
  children,
  imageUrl,
  embedImageUrl,
  imageAlt = "",
  className = "",
  boxClassName = "min-h-[260px]",
  ribbonStrokeWidth = 18,
}: GiftBoxProps) {
  return (
    <div className={`flex flex-col items-center w-full ${className}`}>
      {/* Fixed-size centered ribbon */}
      <GiftRibbon
        strokeWidth={ribbonStrokeWidth}
        className="-mb-[7px] z-10"
      />

      {/* Responsive gift box body */}
      <div
        className={`w-full bg-[#FAFAFA] border-[6px] border-[#1D4B47] rounded-[22px] p-5 shadow-sm transition-all overflow-hidden ${boxClassName}`}
      >
        {children ?? (
          <ProductImage imageUrl={imageUrl} embedImageUrl={embedImageUrl} alt={imageAlt} />
        )}
      </div>
    </div>
  );
}