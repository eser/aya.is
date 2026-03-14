import type React from "react";
import styles from "./center-block.module.css";

type CenterBlockProps = {
  children?: React.ReactNode;
};

function CenterBlock(props: CenterBlockProps) {
  return (
    <div className={styles.center}>{props.children}</div>
  );
}

export { CenterBlock };
export type { CenterBlockProps };
