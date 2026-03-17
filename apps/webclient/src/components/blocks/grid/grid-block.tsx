// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./grid-block.module.css";

type GridBlockProps = {
  cols?: 2 | 3 | 4;
  gap?: "sm" | "md" | "lg";
  minWidth?: string;
  children?: React.ReactNode;
};

function GridBlock(props: GridBlockProps) {
  const useAutoFit = props.minWidth !== undefined;

  const gridStyle: React.CSSProperties = useAutoFit
    ? {
      gridTemplateColumns: `repeat(auto-fit, minmax(${props.minWidth}, 1fr))`,
    }
    : {};

  return (
    <div
      className={styles.grid}
      data-cols={useAutoFit ? undefined : props.cols}
      data-gap={props.gap}
      style={gridStyle}
    >
      {props.children}
    </div>
  );
}

export { GridBlock };
export type { GridBlockProps };
