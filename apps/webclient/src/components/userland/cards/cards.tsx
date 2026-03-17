// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import styles from "./cards.module.css";

export type CardsProps = {
  children?: React.ReactNode;
};

export function Cards(props: CardsProps) {
  return <div className={styles.cards}>{props.children}</div>;
}
