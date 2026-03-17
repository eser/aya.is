// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./gallery-block.module.css";

type GalleryBlockProps = {
  cols?: 2 | 3 | 4;
  gap?: "sm" | "md" | "lg";
  children?: React.ReactNode;
};

function GalleryBlock(props: GalleryBlockProps) {
  return (
    <div
      className={styles.gallery}
      data-cols={props.cols ?? 3}
      data-gap={props.gap ?? "md"}
    >
      {props.children}
    </div>
  );
}

export { GalleryBlock };
export type { GalleryBlockProps };
