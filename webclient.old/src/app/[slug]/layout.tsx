import * as React from "react";
import Image from "next/image";
import { notFound } from "next/navigation";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import { Icons } from "@/shared/components/icons.tsx";

import styles from "./layout.module.css";

function findIcon(kind: string) {
  switch (kind) {
    case "github":
      return Icons.github;
    case "twitter":
    case "x":
      return Icons.twitter;
    case "linkedin":
      return Icons.linkedin;
    case "instagram":
      return Icons.instagram;
    case "youtube":
      return Icons.youtube;
    // case "bsky":
    //   return Icons.bsky;
    default:
      return Icons.link;
  }
}

type LayoutProps = {
  children: React.ReactNode;
  params: Promise<{
    slug: string;
  }>;
};

async function Layout(props: LayoutProps) {
  const params = await props.params;

  const { t, locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);

  if (profileData === null) {
    notFound();
  }

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="grid grid-cols-1 md:grid-cols-[280px_1fr] gap-8 items-start">
        <aside className={styles.bio}>
          {profileData.profile_picture_uri && (
            <div className={styles["profile-picture"]}>
              <Image
                src={profileData.profile_picture_uri}
                alt={`${profileData.title}'s profile picture`}
                width={280}
                height={280}
              />
            </div>
          )}

          <div className={styles.hero}>
            <h1 className={styles.title}>{profileData.title}</h1>

            <div className={styles.subtitle}>
              {profileData.slug}
              {profileData.pronouns && ` · ${profileData.pronouns}`}
            </div>

            {profileData.links?.length > 0 && (
              <div className={styles.links}>
                {profileData.links?.map((link) => {
                  const Icon = findIcon(link.kind);

                  return (
                    <a key={link.id} href={link.uri ?? ""} title={link.title}>
                      <Icon />
                    </a>
                  );
                })}
              </div>
            )}

            <p className={styles.description}>{profileData.description}</p>
          </div>

          <nav className={styles.nav}>
            <ul>
              <li>
                <SiteLink href={`/${params.slug}`}>{t("Layout", "Profile")}</SiteLink>
              </li>

              {profileData.kind === "individual" && (
                <li>
                  <SiteLink href={`/${params.slug}/contributions`}>{t("Layout", "Contributions")}</SiteLink>
                </li>
              )}

              {(profileData.kind === "organization" || profileData.kind === "product") && (
                <li>
                  <SiteLink href={`/${params.slug}/members`}>{t("Layout", "Members")}</SiteLink>
                </li>
              )}

              {profileData.pages?.map((page) => (
                <li key={page.slug}>
                  <SiteLink href={`/${params.slug}/${page.slug}`}>{page.title}</SiteLink>
                </li>
              ))}
            </ul>
          </nav>
        </aside>

        <main>{props.children}</main>
      </div>
    </section>
  );
}

export { Layout as default };
