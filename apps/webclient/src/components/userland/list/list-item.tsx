import { formatMonthShort } from "@/lib/date";
import styles from "./list-item.module.css";

export type ListItemData = {
  date: Date;
  thumbnail: string;
  title: string;
  description: string;
  url: string;
  host: string;
  participants: string[];
};

export type ListItemProps = {
  item: ListItemData;
  locale?: string;
};

export function ListItem(props: ListItemProps) {
  const date = new Date(props.item.date);
  const locale = props.locale ?? "en";

  return (
    <div className={styles.listItem}>
      <div className="sm:w-20 md:w-32">
        <div className={styles.date}>
          <span className={styles.day}>{date.getDate()}</span>
          <span className={styles.monthYear}>
            {formatMonthShort(date, locale)} {date.getFullYear()}
          </span>
        </div>
      </div>
      <div className="flex-1 flex flex-col">
        <h4 className={styles.title}>
          <a
            href={props.item.url}
            title={props.item.title}
            target="_blank"
            rel="noreferrer"
          >
            {props.item.title}
          </a>
        </h4>
        <div className="mb-1">
          <strong>{props.item.host}</strong>
          {props.item.participants.length > 0 && (
            <>
              {" "}<em>with</em>{" "}
              {props.item.participants.join(", ")}
            </>
          )}
        </div>
        <p className={styles.desc}>{props.item.description}</p>
      </div>
    </div>
  );
}
