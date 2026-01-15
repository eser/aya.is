import { LocaleLink } from "@/components/locale-link";
import type { Profile } from "@/modules/backend/types";
import styles from "./profile-card.module.css";

export type ProfileCardProps = {
  profile: Profile;
};

export function ProfileCard(props: ProfileCardProps) {
  return (
    <LocaleLink
      role="card"
      to={`/${props.profile.slug}`}
      className={styles.cardLink}
    >
      <div className={styles.profileCard}>
        {props.profile.profile_picture_uri !== null &&
          props.profile.profile_picture_uri !== undefined && (
          <div className={styles.avatarContainer}>
            <img
              src={props.profile.profile_picture_uri}
              alt={`${props.profile.title}'s picture`}
              width={80}
              height={80}
              className={styles.avatar}
            />
          </div>
        )}
        <div className={styles.info}>
          <h3 className={styles.title}>{props.profile.title}</h3>
          {props.profile.description !== null &&
            props.profile.description !== undefined && <p className={styles.description}>{props.profile.description}
          </p>}
        </div>
      </div>
    </LocaleLink>
  );
}
