// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./columns-block.module.css";

interface ColumnsBlockProps {
  cols?: 2 | 3;
  gap?: "sm" | "md" | "lg";
  children?: React.ReactNode;
}

function ColumnsBlock(props: ColumnsBlockProps) {
  const cols = props.cols ?? 2;
  const gap = props.gap ?? "md";

  return (
    <div className={styles.columns} data-cols={cols} data-gap={gap}>
      {props.children}
    </div>
  );
}

export { ColumnsBlock };
export type { ColumnsBlockProps };
