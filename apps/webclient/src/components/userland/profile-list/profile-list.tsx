// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import styles from "./profile-list.module.css";

export type ProfileListProps = {
  children: React.ReactNode;
};

export function ProfileList(props: ProfileListProps) {
  return <div className={styles.grid}>{props.children}</div>;
}
