import styles from "./profile-list.module.css";

export type ProfileListProps = {
  children: React.ReactNode;
};

export function ProfileList(props: ProfileListProps) {
  return <div className={styles.grid}>{props.children}</div>;
}
