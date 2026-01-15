import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { ProfileListItem } from "@/shared/components/userland/profile-list/profile-list-item.tsx";
import { ProfileList } from "@/shared/components/userland/profile-list/profile-list.tsx";
import styles from "./page.module.css";

type IndexPageProps = {
  params: Promise<{
    slug: string;
  }>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export async function generateMetadata(props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  const params = await props.params;

  const { t, locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);
  if (profileData === null) {
    notFound();
  }

  if (profileData.kind !== "organization" && profileData.kind !== "product") {
    notFound();
  }

  return {
    title: `${profileData.title} - ${t("Layout", "Members")}`,
    description: profileData.description,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { t, locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);
  if (profileData === null) {
    notFound();
  }

  if (profileData.kind !== "organization" && profileData.kind !== "product") {
    notFound();
  }

  const membershipData = await backend.getProfileMembers(locale.code, params.slug);

  if (membershipData === null) {
    notFound();
  }

  return (
    <article className="content space-y-6">
      <div className="flex flex-col gap-2">
        <h2>{t("Layout", "Members")}</h2>
        <h3 className="text-muted-foreground">
          {t("Members", "Individuals and organizations that are members of this profile.")}
        </h3>
      </div>

      {membershipData.length > 0
        ? (
          <ProfileList>
            {membershipData.map((membership) => (
              <ProfileListItem
                key={membership.id}
                focus="member"
                profile={membership.profile}
                membership={membership}
              />
            ))}
          </ProfileList>
        )
        : (
          <div className={styles.emptyState}>
            <p>{t("Members", "No members found.")}</p>
          </div>
        )}
    </article>
  );
}

export { IndexPage as default };
