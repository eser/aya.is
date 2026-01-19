import * as React from "react";
import { cn } from "@/lib/utils";
import { Logo as LogoSvg } from "@/components/icons";

type LogoProps = React.HTMLAttributes<HTMLDivElement>;

export function Logo(props: LogoProps) {
  const { className, ...restProps } = props;

  return (
    <div
      className={cn("flex gap-2 items-center whitespace-nowrap", className)}
      {...restProps}
    >
      <LogoSvg className="w-6 h-6" />
    </div>
  );
}
