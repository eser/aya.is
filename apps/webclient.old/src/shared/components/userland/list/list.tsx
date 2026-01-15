import * as React from "react";

import styles from "./list.module.css";

export type ListProps = {
  children: React.ReactNode;
};

export function List(props: ListProps) {
  return <section className={styles.list}>{props.children}</section>;
}
