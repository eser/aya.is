// Profile members page
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";

export const Route = createFileRoute("/$locale/$slug/members")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const members = await backend.getProfileMembers(locale, slug);
    return { members };
  },
  component: MembersPage,
});

function MembersPage() {
  const { members } = Route.useLoaderData();
  const { t } = useTranslation();

  return (
    <div className="content">
      <h2>{t("Layout.Members")}</h2>
      <p className="text-muted-foreground mb-4">
        {t("Members.Individuals and organizations that are members of this profile.")}
      </p>

      {members && members.length > 0 ? (
        <div className="space-y-4">
          {/* Members will be rendered here */}
        </div>
      ) : (
        <p className="text-muted-foreground">
          {t("Members.No members found.")}
        </p>
      )}
    </div>
  );
}
