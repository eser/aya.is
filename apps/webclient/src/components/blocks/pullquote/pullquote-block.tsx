import type React from "react";
import styles from "./pullquote-block.module.css";

type PullquoteBlockProps = {
  citation?: string;
  children?: React.ReactNode;
};

function PullquoteBlock(props: PullquoteBlockProps) {
  return (
    <figure className={styles.pullquote}>
      <blockquote className={styles.text}>{props.children}</blockquote>
      {props.citation !== undefined && props.citation !== "" && (
        <figcaption className={styles.citation}>— {props.citation}</figcaption>
      )}
    </figure>
  );
}

export { PullquoteBlock };
export type { PullquoteBlockProps };
