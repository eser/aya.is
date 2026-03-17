// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import { ChevronRight } from "lucide-react";
import styles from "./details-block.module.css";

interface DetailsBlockProps {
  summary: string;
  open?: boolean;
  children?: React.ReactNode;
}

function DetailsBlock(props: DetailsBlockProps) {
  return (
    <details className={styles.details} open={props.open === true ? true : undefined}>
      <summary>
        <ChevronRight className={styles.chevron} />
        {props.summary}
      </summary>
      <div className={styles.content}>{props.children}</div>
    </details>
  );
}

export { DetailsBlock };
export type { DetailsBlockProps };
