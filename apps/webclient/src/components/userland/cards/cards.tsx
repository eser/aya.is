import styles from "./cards.module.css";

export type CardsProps = {
  children?: React.ReactNode;
};

export function Cards(props: CardsProps) {
  return <div className={styles.cards}>{props.children}</div>;
}
