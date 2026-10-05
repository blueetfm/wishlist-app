import * as React from "react";

import { cn } from "@/lib/utils";
import "./stylized-button.css";

type StylizedButtonVariant = "outline" | "solid";
type StylizedButtonSize = "sm" | "default" | "lg";

interface StylizedButtonProps extends React.ComponentProps<"button"> {
  variant?: StylizedButtonVariant;
  size?: StylizedButtonSize;
}

export function StylizedButton({
  className,
  variant = "outline",
  size = "default",
  ...props
}: StylizedButtonProps) {
  return (
    <button
      className={cn(
        "stylized-button",
        `stylized-button--${variant}`,
        `stylized-button--${size}`,
        className
      )}
      {...props}
    />
  );
}

export default StylizedButton;
