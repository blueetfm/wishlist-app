"use client"

import React from "react";

import Giftbox from "../gift-box/gift-box"
import "./wishlist-item-summary.css"

interface WishlistItemSummaryProps {
    name: string
    price: number
    imageUrl?: string
    imageAlt?: string
}

export default function WishlistItemSummary({
    name,
    price,
    imageUrl,
    imageAlt
}: WishlistItemSummaryProps) {
    return (
        <div className="wishlist-item-summary">
            <Giftbox imageUrl={imageUrl} />

            <div className="wishlist-item-summary__overlay">
                <h3 className="wishlist-item-summary__title">{name}</h3>
                <p className="wishlist-item-summary__price">${price.toFixed(2)}</p>
            </div>
            
        </div>
    );
}