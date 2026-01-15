import * as React from "react";
import NextImage from "next/image";
import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import styles from "./profile-card.module.css";

const Image = NextImage;

// Placeholder - this should be refined based on actual backend types
export type Profile = {
  slug: string;
  title: string;
  description?: string | null;
  profile_picture_uri?: string | null;
  kind?: "individual" | "organization" | string;
};

export type ProfileCardProps = {
  profile: Profile;
};

export function ProfileCard(props: ProfileCardProps) {
  const { profile } = props;

  return (
    <SiteLink role="card" href={`/${profile.slug}`} className={styles.cardLink}>
      <div className={styles.profileCard}>
        {profile.profile_picture_uri && (
          <div className={styles.avatarContainer}>
            <Image
              src={profile.profile_picture_uri}
              alt={`${profile.title}'s picture`}
              width={80}
              height={80}
              className={styles.avatar}
            />
          </div>
        )}
        <div className={styles.info}>
          <h3 className={styles.title}>{profile.title}</h3>
          {profile.description && <p className={styles.description}>{profile.description}</p>}
        </div>
      </div>
    </SiteLink>
  );
}
