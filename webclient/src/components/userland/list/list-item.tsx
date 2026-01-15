import { formatDate } from "@/lib/date";
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
            {formatDate(date, locale)}
          </span>
        </div>
      </div>
      <div className={`flex-1 ${props.item.thumbnail !== "" ? "flex flex-row gap-4" : ""}`}>
        {props.item.thumbnail !== "" && (
          <img className={styles.media} src={props.item.thumbnail} alt="media thumbnail" />
        )}
        <div className="flex flex-col mb-2">
          <h4 className={styles.title}>
            <a href={props.item.url} title={props.item.title} target="_blank" rel="noreferrer">
              {props.item.title}
            </a>
          </h4>
          <div>
            <strong>{props.item.host}</strong>
            {props.item.participants.length > 0 && <em>{" ve "}</em>}
            {props.item.participants.length > 0 && props.item.participants.join(", ")}
          </div>
        </div>
        <p className={styles.desc}>{props.item.description}</p>
      </div>
    </div>
  );
}
