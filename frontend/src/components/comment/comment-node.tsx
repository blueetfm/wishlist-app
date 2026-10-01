import React from "react";

export default interface CommentNode {
    id: string;

    author: string;
    content: string;
    timestamp: string;
    children: CommentNode[];
}