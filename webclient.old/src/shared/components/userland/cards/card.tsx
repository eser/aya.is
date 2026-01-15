import * as React from "react";

import Image from "next/image";

import { Conditional } from "@/shared/components/conditional.tsx";

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
        <Conditional test={props.category !== undefined}>
          <div className={styles.tags}>{props.category}</div>
        </Conditional>
        <Conditional test={props.imageUri !== undefined && props.imageUri !== null}>
          <div className={styles.image}>
            <Image src={props.imageUri as string} width={220} height={220} alt={props.title} loading="lazy" />
          </div>
        </Conditional>
        <h5 className={styles.title}>{props.title}</h5>
        <p className={styles.description}>{props.description}</p>
      </div>
    </a>
  );
}
