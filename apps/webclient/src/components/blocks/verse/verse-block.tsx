import type React from "react";
import styles from "./verse-block.module.css";

type VerseBlockProps = {
  children?: React.ReactNode;
};

function VerseBlock(props: VerseBlockProps) {
  return (
    <div className={styles.verse}>{props.children}</div>
  );
}

export { VerseBlock };
export type { VerseBlockProps };
