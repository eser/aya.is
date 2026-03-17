// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./button-block.module.css";

interface ButtonBlockProps {
  href: string;
  variant?: "default" | "outline" | "secondary";
  size?: "sm" | "default" | "lg";
  children?: React.ReactNode;
}

const SAFE_HREF_PREFIXES = ["http://", "https://", "/", "#"];

function isSafeHref(href: string): boolean {
  return SAFE_HREF_PREFIXES.some((prefix) => href.startsWith(prefix));
}

function isExternalHref(href: string): boolean {
  return href.startsWith("http://") || href.startsWith("https://");
}

function ButtonBlock(props: ButtonBlockProps) {
  const variant = props.variant ?? "default";
  const size = props.size ?? "default";

  if (!isSafeHref(props.href)) {
    return <span className={styles.plainText}>{props.children}</span>;
  }

  const externalProps = isExternalHref(props.href) ? { target: "_blank" as const, rel: "noopener noreferrer" } : {};

  return (
    <a
      href={props.href}
      className={styles.button}
      data-variant={variant}
      data-size={size}
      {...externalProps}
    >
      {props.children}
    </a>
  );
}

export { ButtonBlock };
export type { ButtonBlockProps };
