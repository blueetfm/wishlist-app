import React from "react";
import "./gift-ribbon.css";

interface GiftRibbonProps {
  /** Controls the thickness of the ribbon lines (default: 18) */
  strokeWidth?: number;
  /** Stroke color (default: #1D4B47 dark teal) */
  color?: string;
  /** Optional class overrides for positioning or sizing */
  className?: string;
}

export function GiftRibbon({
  strokeWidth = 18,
  color = "#1D4B47",
  className = "",
}: GiftRibbonProps) {
  return (
    <div
      className={`gift-ribbon ${className}`.trim()}
    >
      <svg
        viewBox="0 0 210 70"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="gift-ribbon__svg"
      >
        {/* Left Ribbon Tail */}
        <path
          d="M98 62 C68 62, 44 52, 20 36"
          stroke={color}
          strokeWidth={strokeWidth}
          strokeLinecap="butt"
        />
        {/* Right Ribbon Tail */}
        <path
          d="M112 62 C142 62, 166 52, 190 36"
          stroke={color}
          strokeWidth={strokeWidth}
          strokeLinecap="butt"
        />
        {/* Left Bow Loop */}
        <path
          d="M102 62 C58 56, 42 28, 58 14 C74 0, 92 26, 102 62 Z"
          stroke={color}
          strokeWidth={strokeWidth}
          strokeLinejoin="round"
        />
        {/* Right Bow Loop */}
        <path
          d="M108 62 C152 56, 168 28, 152 14 C136 0, 118 26, 108 62 Z"
          stroke={color}
          strokeWidth={strokeWidth}
          strokeLinejoin="round"
        />
      </svg>
    </div>
  );
}

export default GiftRibbon;