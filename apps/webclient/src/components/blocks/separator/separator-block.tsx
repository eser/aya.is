// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import styles from "./separator-block.module.css";

interface SeparatorBlockProps {
  variant?: "line" | "dots" | "space";
}

function SeparatorBlock(props: SeparatorBlockProps) {
  const variant = props.variant ?? "line";

  if (variant === "line") {
    return (
      <hr
        className={styles.separator}
        data-variant="line"
        role="separator"
        aria-hidden="true"
      />
    );
  }

  if (variant === "dots") {
    return (
      <div
        className={styles.separator}
        data-variant="dots"
        role="separator"
        aria-hidden="true"
      >
        · · ·
      </div>
    );
  }

  return (
    <div
      className={styles.separator}
      data-variant="space"
      role="separator"
      aria-hidden="true"
    />
  );
}

export { SeparatorBlock };
export type { SeparatorBlockProps };
