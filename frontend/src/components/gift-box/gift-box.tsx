import React from "react";
import { GiftRibbon } from "./gift-ribbon";
import { ProductImage } from "./product-image";
import "./gift-box.css";

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
  boxClassName = "gift-box__box--default-height",
  ribbonStrokeWidth = 18,
}: GiftBoxProps) {
  return (
    <div className={`gift-box ${className}`.trim()}>
      {/* Fixed-size centered ribbon */}
      <GiftRibbon
        strokeWidth={ribbonStrokeWidth}
        className="gift-box__ribbon"
      />

      {/* Responsive gift box body */}
      <div
        className={`gift-box__box ${boxClassName}`.trim()}
      >
        {children ?? (
          <ProductImage imageUrl={imageUrl} embedImageUrl={embedImageUrl} alt={imageAlt} />
        )}
      </div>
    </div>
  );
}