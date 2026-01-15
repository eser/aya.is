"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

type ToggleGroupContextValue = {
  value: string;
  onValueChange: (value: string) => void;
  variant?: "default" | "outline";
};

const ToggleGroupContext = React.createContext<ToggleGroupContextValue | null>(
  null
);

export type ToggleGroupProps = {
  type: "single";
  value: string;
  onValueChange: (value: string) => void;
  variant?: "default" | "outline";
  className?: string;
  children: React.ReactNode;
  "aria-label"?: string;
  id?: string;
};

function ToggleGroup({
  value,
  onValueChange,
  variant = "default",
  className,
  children,
  ...props
}: ToggleGroupProps) {
  return (
    <div
      data-slot="toggle-group"
      data-variant={variant}
      role="group"
      className={cn(
        "inline-flex rounded-md",
        variant === "outline" && "shadow-xs",
        className
      )}
      {...props}
    >
      <ToggleGroupContext.Provider value={{ value, onValueChange, variant }}>
        {children}
      </ToggleGroupContext.Provider>
    </div>
  );
}

export type ToggleGroupItemProps = {
  value: string;
  className?: string;
  children: React.ReactNode;
  "aria-label"?: string;
};

function ToggleGroupItem({
  value,
  className,
  children,
  ...props
}: ToggleGroupItemProps) {
  const context = React.useContext(ToggleGroupContext);
  if (!context) {
    throw new Error("ToggleGroupItem must be used within a ToggleGroup");
  }

  const isSelected = context.value === value;
  const variant = context.variant;

  return (
    <button
      type="button"
      data-slot="toggle-group-item"
      data-state={isSelected ? "on" : "off"}
      data-variant={variant}
      role="radio"
      aria-checked={isSelected}
      onClick={() => context.onValueChange(value)}
      className={cn(
        "inline-flex items-center justify-center text-sm font-medium transition-colors",
        "h-9 px-3 min-w-9",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:pointer-events-none disabled:opacity-50",
        "first:rounded-l-md last:rounded-r-md",
        variant === "outline" && [
          "border border-input bg-transparent",
          "border-l-0 first:border-l",
          isSelected
            ? "bg-accent text-accent-foreground"
            : "hover:bg-accent hover:text-accent-foreground",
        ],
        variant === "default" && [
          "bg-transparent",
          isSelected
            ? "bg-accent text-accent-foreground"
            : "hover:bg-accent hover:text-accent-foreground",
        ],
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}

export { ToggleGroup, ToggleGroupItem };
