import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const stylizedButtonVariants = cva(
  "inline-flex items-center justify-center rounded-xl border-2 border-[#214642] font-semibold transition-all disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        outline:
          "bg-transparent text-[#214642] shadow-[3px_3px_0px_#214642] hover:translate-y-[1px] hover:shadow-[2px_2px_0px_#214642]",
        solid:
          "bg-[#214642] text-[#FAF9F5] shadow-[4px_4px_0px_rgba(33,70,66,0.3)] hover:bg-[#2a5752]",
      },
      size: {
        sm: "px-6 py-2.5 text-sm uppercase tracking-widest",
        default: "py-3 px-6",
        lg: "w-full py-4 px-6 text-lg",
      },
    },
    defaultVariants: {
      variant: "outline",
      size: "default",
    },
  }
);

interface StylizedButtonProps
  extends React.ComponentProps<"button">,
    VariantProps<typeof stylizedButtonVariants> {}

/** The dark-teal, drop-shadow "pop" button used throughout the marketing page. */
export function StylizedButton({
  className,
  variant,
  size,
  ...props
}: StylizedButtonProps) {
  return (
    <button
      className={cn(stylizedButtonVariants({ variant, size, className }))}
      {...props}
    />
  );
}

export { stylizedButtonVariants };
export default StylizedButton;
