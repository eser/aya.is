// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type React from "react";
import styles from "./verse-block.module.css";

type VerseBlockProps = {
  children?: React.ReactNode;
};

function VerseBlock(props: VerseBlockProps) {
  return <div className={styles.verse}>{props.children}</div>;
}

export { VerseBlock };
export type { VerseBlockProps };
