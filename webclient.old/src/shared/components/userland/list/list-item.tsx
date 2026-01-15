import * as React from "react";
import Image from "next/image";

import { getMonthName } from "@/shared/lib/date.ts";

import styles from "./list-item.module.css";

export type Item = {
  date: Date;
  thumbnail: string;
  title: string;
  description: string;
  url: string;
  host: string;
  participants: string[];
};

export type ListItemProps = {
  item: Item;
};

export function ListItem(props: ListItemProps) {
  return (
    <div className={styles["list-item"]}>
      <div className="sm:w-20 md:w-32">
        <div className={styles.date}>
          <span className={styles.day}>{new Date(props.item.date).getDate()}</span>
          <span className={styles["month-n-year"]}>
            {`${getMonthName(props.item.date)} ${new Date(props.item.date).getFullYear()}`}
          </span>
        </div>
      </div>
      <div className={`flex-1 ${props.item.thumbnail ? "flex flex-row gap-4" : ""}`}>
        {props.item.thumbnail && <Image className={styles.media} src={props.item.thumbnail} alt="media thumbnail" />}
        <div className="flex flex-col mb-2">
          <h4 className={`${styles.title} dark:text-theme-dark-primary`}>
            <a href={props.item.url} title={props.item.title} target="_blank">
              {props.item.title}
            </a>
          </h4>
          <div>
            <strong>{props.item.host}</strong>
            {props.item.participants.length > 0 && <em>{` ve `}</em>}
            {props.item.participants.length > 0 && `${props.item.participants.join(", ")}`}
          </div>
        </div>
        <p className={styles.desc}>{props.item.description}</p>
      </div>
    </div>
  );
}
