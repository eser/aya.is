// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./stack-block.module.css";

type StackBlockProps = {
  gap?: "sm" | "md" | "lg";
  align?: "start" | "center" | "end" | "stretch";
  children?: React.ReactNode;
};

function StackBlock(props: StackBlockProps) {
  return (
    <div
      className={styles.stack}
      data-gap={props.gap}
      data-align={props.align}
    >
      {props.children}
    </div>
  );
}

export { StackBlock };
export type { StackBlockProps };
