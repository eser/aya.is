// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./row-block.module.css";

type RowBlockProps = {
  gap?: "sm" | "md" | "lg";
  align?: "start" | "center" | "end" | "stretch";
  justify?: "start" | "center" | "end" | "between" | "around";
  wrap?: boolean;
  children?: React.ReactNode;
};

function RowBlock(props: RowBlockProps) {
  return (
    <div
      className={styles.row}
      data-gap={props.gap}
      data-align={props.align}
      data-justify={props.justify}
      data-wrap={props.wrap !== undefined ? String(props.wrap) : undefined}
    >
      {props.children}
    </div>
  );
}

export { RowBlock };
export type { RowBlockProps };
