import * as React from "react";

import styles from "./profile-list.module.css";

interface ProfileListProps {
  children: React.ReactNode;
}

export function ProfileList(props: ProfileListProps) {
  return (
    <div className={styles.grid}>
      {props.children}
    </div>
  );
}
