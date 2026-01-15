import * as React from "react";

import { siteConfig } from "@/shared/config.ts";
import { cn } from "@/shared/lib/cn.ts";
import { Icons } from "@/shared/components/icons.tsx";

import styles from "./logo.module.css";
import { EasterEgg } from "./easter-egg.tsx";

export function Logo(props: React.HTMLAttributes<HTMLElement>) {
  return (
    <div className={cn(styles.container, props.className)}>
      <Icons.logo className="w-6 h-6" />

      <span className="text-base font-bold site-name">{siteConfig.name}</span>

      <EasterEgg />
    </div>
  );
}
