"use client"

import React from "react"
import { useIsMobile } from "../../lib/is-mobile"
import CommentNode from "../comment/comment-node"
import "./item-description.css"

interface ItemDescriptionProps {
    title: string
    description?: string
    price: number
    children?: React.ReactNode
    // comments: CommentNode[]
}

export default function ItemDescription({
    title,
    description,
    price,
    children
}: ItemDescriptionProps) {
    return (
        <div className="item-description">
            <h3 className="item-description__title">
              {title}
            </h3>
            {description && <p className="item-description__description">{description}</p>}
            <p className="item-description__price">${price.toFixed(2)}</p>
            {children}
        </div>
    )
}