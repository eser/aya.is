// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./group-block.module.css";

type GroupBlockProps = {
  padding?: "sm" | "md" | "lg";
  background?: "muted" | "card" | "accent";
  border?: boolean;
  rounded?: boolean;
  children?: React.ReactNode;
};

function GroupBlock(props: GroupBlockProps) {
  return (
    <div
      className={styles.group}
      data-padding={props.padding}
      data-background={props.background}
      data-border={props.border === true ? "true" : undefined}
      data-rounded={props.rounded === true ? "true" : undefined}
    >
      {props.children}
    </div>
  );
}

export { GroupBlock };
export type { GroupBlockProps };
