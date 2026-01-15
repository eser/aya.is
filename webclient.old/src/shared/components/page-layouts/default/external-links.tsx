import * as React from "react";

import NextLinkOrig from "next/link";
// More robust workaround for Deno/Next.js Link typing issues
const Link = NextLinkOrig as unknown as typeof import("next/link").default;

import { siteConfig } from "@/shared/config.ts";
import { cn } from "@/shared/lib/cn.ts";
import { Icons } from "@/shared/components/icons.tsx";
import { Button } from "@/shared/components/ui/button.tsx";

import styles from "./external-links.module.css";

type NavItem = {
  key: string;
  title: string;
  href?: string;
  icon: React.ReactNode;
  disabled?: boolean;
};

const navItems: NavItem[] = [
  {
    key: "github",
    title: "GitHub",
    icon: <Icons.github className={styles.icon} />,
    href: siteConfig.links.github,
  },
  {
    key: "x",
    title: "X",
    icon: <Icons.twitter className={styles.icon} />,
    href: siteConfig.links.x,
  },
  {
    key: "instagram",
    title: "Instagram",
    icon: <Icons.instagram className={styles.icon} />,
    href: siteConfig.links.instagram,
  },
];

export function ExternalLinks() {
  return (
    <div className={styles.container}>
      {navItems.map(
        (item, _index) =>
          item.href && (
            <Button key={item.key} variant="ghost" size="sm" asChild>
              <Link href={item.href} target="_blank" rel="noreferrer" className={cn(item.disabled && styles.disabled)}>
                {item.icon}
                <span className={styles.description}>{item.title}</span>
              </Link>
            </Button>
          ),
      )}
    </div>
  );
}
