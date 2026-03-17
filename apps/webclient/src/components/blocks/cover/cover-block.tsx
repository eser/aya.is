// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import { cn } from "@/lib/utils";
import styles from "./cover-block.module.css";

type CoverBlockProps = {
  src: string;
  alt?: string;
  overlay?: "dark" | "light" | "none";
  minHeight?: string;
  children?: React.ReactNode;
};

function CoverBlock(props: CoverBlockProps) {
  const overlay = props.overlay ?? "dark";

  return (
    <div
      className={cn(
        styles.cover,
        overlay === "dark" || overlay === undefined ? styles.coverDark : styles.coverLight,
      )}
      style={{
        backgroundImage: `url(${props.src})`,
        minHeight: props.minHeight ?? "300px",
      }}
      role="img"
      aria-label={props.alt ?? ""}
    >
      {overlay !== "none" && <div className={styles.overlay} data-overlay={overlay} />}
      <div className={styles.content}>{props.children}</div>
    </div>
  );
}

export { CoverBlock };
export type { CoverBlockProps };
