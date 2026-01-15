import * as React from "react";

import styles from "./cards.module.css";

type CardsProps = {
  children?: React.ReactNode;
};

export function Cards(props: CardsProps) {
  return <div className={styles.cards}>{props.children}</div>;
}
