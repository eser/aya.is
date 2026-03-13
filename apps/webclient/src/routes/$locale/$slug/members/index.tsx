// Profile members list page
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Send } from "lucide-react";
import { profileMembersQueryOptions, profileQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { MemberCard } from "@/components/userland/member-card/member-card";
import { LocaleLink } from "@/components/locale-link";
import { buildUrl, generateCanonicalLink, generateMetaTags } from "@/lib/seo";
import i18next from "i18next";
import { NotFoundContent } from "./route";

const parentRoute = getRouteApi("/$locale/$slug");

export const Route = createFileRoute("/$locale/$slug/members/")({
  loader: async ({ params, context }) => {
    const { locale, slug } = params;
    const profile = await context.queryClient.ensureQueryData(profileQueryOptions(locale, slug));

    if (profile?.feature_relations === "disabled") {
      return {
        members: null,
        locale,
        slug,
        translatedTitle: "",
        translatedDescription: "",
        notFound: true as const,
      };
    }

    const members = await context.queryClient.ensureQueryData(profileMembersQueryOptions(locale, slug));

    if (members === null) {
      return {
        members: null,
        locale,
        slug,
        translatedTitle: "",
        translatedDescription: "",
        notFound: true as const,
      };
    }

    await i18next.loadLanguages(locale);
    const t = i18next.getFixedT(locale);
    const translatedTitle = `${t("Layout.Members")} - ${profile?.title ?? slug}`;
    const translatedDescription = t(
      "Members.Individuals and organizations that are members of this profile.",
    );

    return {
      members,
      locale,
      slug,
      translatedTitle,
      translatedDescription,
      notFound: false as const,
    };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound) {
      return { meta: [] };
    }

    const { locale, slug, translatedTitle, translatedDescription } = loaderData;

    return {
      meta: generateMetaTags({
        title: translatedTitle,
        description: translatedDescription,
        url: buildUrl(locale, slug, "members"),
        locale,
        type: "website",
      }),
      links: [generateCanonicalLink(buildUrl(locale, slug, "members"))],
    };
  },
  errorComponent: QueryError,
  component: MembersPage,
});

const MEMBER_KINDS = new Set([
  "member",
  "contributor",
  "maintainer",
  "lead",
  "owner",
]);

function MembersPage() {
  const loaderData = Route.useLoaderData();
  const { profile, permissions } = parentRoute.useLoaderData();
  const { t } = useTranslation();
  const params = Route.useParams();

  if (loaderData.notFound || loaderData.members === null || profile === null) {
    return <NotFoundContent />;
  }

  const { members } = loaderData;

  const isMember = permissions?.viewer_membership_kind !== undefined &&
    permissions.viewer_membership_kind !== null &&
    MEMBER_KINDS.has(permissions.viewer_membership_kind);

  const showApplyButton = profile.feature_applications !== "disabled" &&
    !isMember;

  return (
    <>
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-serif text-2xl font-bold text-foreground">
            {t("Layout.Members")}
          </h2>
          <p className="text-muted-foreground">
            {t(
              "Members.Individuals and organizations that are members of this profile.",
            )}
          </p>
        </div>
        {showApplyButton && (
          <LocaleLink
            to={`/${params.slug}/members/apply`}
            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-colors no-underline"
          >
            <Send className="size-4" />
            {t("Applications.Apply to Join")}
          </LocaleLink>
        )}
      </div>

      {members !== null && members.length > 0
        ? (
          <div className="flex flex-col gap-4">
            {members.map((membership) => <MemberCard key={membership.id} membership={membership} />)}
          </div>
        )
        : (
          <p className="text-muted-foreground">
            {t("Members.No members found.")}
          </p>
        )}
    </>
  );
}
