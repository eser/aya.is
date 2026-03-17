// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { Skeleton } from "@/components/ui/skeleton";
import styles from "./profile-listing.module.css";

export function ProfileCardSkeleton() {
  return (
    <div className={styles.skeletonCard}>
      <Skeleton className={styles.skeletonImage} />
      <div className={styles.skeletonInfo}>
        <Skeleton className={styles.skeletonTitle} />
        <Skeleton className={styles.skeletonDesc} />
      </div>
    </div>
  );
}

export type ProfileCardSkeletonGridProps = {
  count?: number;
};

export function ProfileCardSkeletonGrid(props: ProfileCardSkeletonGridProps) {
  const count = props.count ?? 8;

  return (
    <div className={styles.grid}>
      {Array.from({ length: count }, (_, i) => <ProfileCardSkeleton key={i} />)}
    </div>
  );
}
