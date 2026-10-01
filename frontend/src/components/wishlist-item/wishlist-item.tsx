import React from "react";

import Giftbox from "../gift-box/gift-box"
import CommentNode from "../comment/comment-node"

interface WishlistItemProps {
    children? : React.ReactNode
    name: string
    description: string
    price: number
    imageUrl?: string
    // comments: CommentNode[]
}

export default function WishlistItem({
    children,
    name,
    description,
    price,
    imageUrl,
    // comments
}: WishlistItemProps) {
    return (
        <div className="wishlist-item">
            <Giftbox imageUrl={imageUrl} />
            <h3>{name}</h3>
            <p>{description}</p>
            <p>${price.toFixed(2)}</p>
            {children}
        </div>
    );
}