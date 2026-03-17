// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import styles from "./card.module.css";

export type CardProps = {
  category?: string;
  imageUri?: string | null;
  title: string;
  description: string;
  href?: string;
  children?: React.ReactNode;
};

export function Card(props: CardProps) {
  return (
    <a className={styles.card} href={props.href}>
      <div className={styles.inner}>
        {props.category !== undefined && <div className={styles.tags}>{props.category}</div>}
        {props.imageUri !== undefined && props.imageUri !== null && (
          <div className={styles.image}>
            <img
              src={props.imageUri}
              width={220}
              height={220}
              alt={props.title}
              loading="lazy"
            />
          </div>
        )}
        <h5 className={styles.title}>{props.title}</h5>
        <p className={styles.description}>{props.description}</p>
      </div>
    </a>
  );
}
