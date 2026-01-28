// Admin profile general info tab
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { formatDateTimeLong } from "@/lib/date";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

export const Route = createFileRoute("/$locale/admin/profiles/$slug/")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const profile = await backend.getAdminProfile(locale, slug);
    return { profile };
  },
  component: AdminProfileGeneral,
});

function AdminProfileGeneral() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const { profile } = Route.useLoaderData();

  if (profile === null || profile === undefined) {
    return null;
  }

  return (
    <div className="space-y-6">
      {/* Profile Info Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t("Admin.Profile Information")}</CardTitle>
          <CardDescription>
            {t("Admin.Basic profile details")}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>{t("Common.ID")}</Label>
              <Input value={profile.id} readOnly className="font-mono text-sm bg-muted" />
            </div>
            <div className="space-y-2">
              <Label>{t("Common.Slug")}</Label>
              <Input value={profile.slug} readOnly className="font-mono bg-muted" />
            </div>
          </div>

          <div className="space-y-2">
            <Label>{t("Common.Title")}</Label>
            <Input
              value={profile.title}
              readOnly
              className={`bg-muted ${profile.has_translation === false ? "italic text-muted-foreground" : ""}`}
              placeholder={profile.has_translation === false ? t("Admin.no translation found") : undefined}
            />
          </div>

          <div className="space-y-2">
            <Label>{t("Common.Description")}</Label>
            <Textarea
              value={profile.description ?? ""}
              readOnly
              className="bg-muted resize-none"
              rows={3}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>{t("Common.Kind")}</Label>
              <Input value={profile.kind} readOnly className="bg-muted capitalize" />
            </div>
            <div className="space-y-2">
              <Label>{t("Common.Points")}</Label>
              <Input value={profile.points.toLocaleString()} readOnly className="bg-muted" />
            </div>
          </div>

          {profile.pronouns !== null && profile.pronouns !== undefined && (
            <div className="space-y-2">
              <Label>{t("Common.Pronouns")}</Label>
              <Input value={profile.pronouns} readOnly className="bg-muted" />
            </div>
          )}

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>{t("Common.Created")}</Label>
              <Input value={formatDateTimeLong(new Date(profile.created_at), locale)} readOnly className="bg-muted" />
            </div>
            {profile.updated_at !== null && profile.updated_at !== undefined && (
              <div className="space-y-2">
                <Label>{t("Admin.Updated")}</Label>
                <Input value={formatDateTimeLong(new Date(profile.updated_at), locale)} readOnly className="bg-muted" />
              </div>
            )}
          </div>

          <div className="space-y-2">
            <Label>{t("Admin.Translation Status")}</Label>
            <Input
              value={profile.has_translation ? t("Admin.Has translation") : t("Admin.No translation")}
              readOnly
              className={`bg-muted ${profile.has_translation ? "text-green-600" : "text-orange-600"}`}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
