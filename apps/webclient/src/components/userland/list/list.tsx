// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import styles from "./list.module.css";

export type ListProps = {
  children: React.ReactNode;
};

export function List(props: ListProps) {
  return <section className={styles.list}>{props.children}</section>;
}
