// Profile members layout with Members/Referrals tab navigation
import { createFileRoute, getRouteApi, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { LocaleLink } from "@/components/locale-link";
import { ProfileSidebarLayout } from "@/components/profile-sidebar-layout";
import { ChildNotFound } from "../route";

const parentRoute = getRouteApi("/$locale/$slug");

const MEMBER_PLUS_KINDS = new Set([
  "member",
  "contributor",
  "maintainer",
  "lead",
  "owner",
]);

export const Route = createFileRoute("/$locale/$slug/members")({
  component: MembersLayout,
  notFoundComponent: MembersChildNotFound,
});

function MembersChildNotFound() {
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Page not found")}</h2>
      <p className="text-muted-foreground">
        {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
      </p>
    </div>
  );
}

function MembersLayout() {
  const { profile, permissions } = parentRoute.useLoaderData();
  const { t } = useTranslation();
  const params = Route.useParams();

  if (
    profile === null || profile.feature_relations === "disabled"
  ) {
    return <ChildNotFound />;
  }

  const isMemberPlus = permissions?.viewer_membership_kind !== undefined &&
    MEMBER_PLUS_KINDS.has(permissions.viewer_membership_kind);

  return (
    <ProfileSidebarLayout
      profile={profile}
      slug={params.slug}
      locale={params.locale}
      viewerMembershipKind={permissions?.viewer_membership_kind}
    >
      <div className="space-y-6">
        {isMemberPlus && (
          <nav className="flex gap-4 border-b">
            <LocaleLink
              to={`/${params.slug}/members`}
              activeOptions={{ exact: true }}
              className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
              activeProps={{
                className:
                  "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
              }}
            >
              {t("Layout.Members")}
            </LocaleLink>
            <LocaleLink
              to={`/${params.slug}/members/referrals`}
              className="relative pb-2 text-sm font-medium text-muted-foreground hover:text-foreground"
              activeProps={{
                className:
                  "relative pb-2 text-sm font-medium text-foreground after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-foreground",
              }}
            >
              {t("Layout.Referrals")}
            </LocaleLink>
          </nav>
        )}

        <Outlet />
      </div>
    </ProfileSidebarLayout>
  );
}
