// Profile members page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { MemberCard } from "@/components/userland/member-card/member-card";
import { buildUrl, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";

const parentRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/members")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const members = await backend.getProfileMembers(locale, slug);
    const profile = await backend.getProfile(locale, slug);

    // Ensure locale translations are loaded before translating
    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Members")} - ${profile?.title ?? slug}`;
    const translatedDescription = t("Members.Individuals and organizations that are members of this profile.");

    return {
      members,
      locale,
      slug,
      profileTitle: profile?.title ?? slug,
      translatedTitle,
      translatedDescription,
    };
  },
  head: ({ loaderData }) => {
    const { locale, slug, translatedTitle, translatedDescription } = loaderData;
    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, slug, "members"),
        locale,
        type: "website",
      }),
    };
  },
  component: MembersPage,
});

function MembersPage() {
  const { members, locale, slug } = Route.useLoaderData();
  const { profile } = parentRoute.useLoaderData();
  const { t } = useTranslation();

  if (profile === null) {
    return null;
  }

  return (
    <ProfileSidebarLayout profile={profile} slug={slug} locale={locale}>
      <div className="space-y-6">
        <div>
          <h2 className="font-serif text-2xl font-bold text-foreground">{t("Layout.Members")}</h2>
          <p className="text-muted-foreground">
            {t(
              "Members.Individuals and organizations that are members of this profile.",
            )}
          </p>
        </div>

        {members !== null && members.length > 0
          ? (
            <div className="flex flex-col gap-4">
              {members.map((membership) => (
                <MemberCard key={membership.id} membership={membership} />
              ))}
            </div>
          )
          : (
            <p className="text-muted-foreground">
              {t("Members.No members found.")}
            </p>
          )}
      </div>
    </ProfileSidebarLayout>
  );
}
