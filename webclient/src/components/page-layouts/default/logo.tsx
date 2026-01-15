import * as React from "react";
import { siteConfig } from "@/config";
import { cn } from "@/lib/utils";
import { Icons } from "@/components/icons";

interface LogoProps extends React.HTMLAttributes<HTMLDivElement> {}

export function Logo({ className, ...props }: LogoProps) {
  return (
    <div className={cn("flex gap-2 items-center whitespace-nowrap", className)} {...props}>
      <Icons.logo className="w-6 h-6" />
    </div>
  );
}
