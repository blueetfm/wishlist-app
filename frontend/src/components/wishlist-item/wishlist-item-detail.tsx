"use client"

import React from "react";
import { useIsMobile } from "../../lib/is-mobile"
import Giftbox from "../gift-box/gift-box"
import ItemDescription from "./item-description"
import CommentNode from "../comment/comment-node"
import "./wishlist-item-detail.css"

interface WishlistItemDetailProps {
    children? : React.ReactNode
    name: string
    description?: string
    price: number
    imageUrl?: string
    imageAlt?: string
    // comments: CommentNode[]
}

export default function WishlistItemDetail({
    children,
    name,
    description,
    price,
    imageUrl,
    imageAlt
    // comments
}: WishlistItemDetailProps) {
    const isMobile = useIsMobile()

    return (
        <div className="wishlist-item-detail">
            { isMobile ? (
                <div className="wishlist-item-detail__column">  
                {/* if on mobile, giftbox on top, container at the bottom */}
                    <div className="wishlist-item-detail__giftbox">
                        <Giftbox imageUrl={imageUrl} />
                    </div>
                    
                    <ItemDescription title={name} description={description} price={price}>
                        {children}
                    </ItemDescription>
                </div>
            ): (
                <div className="wishlist-item-detail__row">
                    {/* Container on the left */}
                    <div className="wishlist-item-detail__giftbox">
                        <Giftbox imageUrl={imageUrl} />
                    </div>
                    {/* Container on the right, capped to the giftbox's height */}
                    <ItemDescription title={name} description={description} price={price}>
                        {children}
                    </ItemDescription>
                </div>
            )}
            
        </div>
    );
}