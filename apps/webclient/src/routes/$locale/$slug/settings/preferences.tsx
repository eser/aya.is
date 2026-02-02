// Profile settings preferences
import * as React from "react";
import { createFileRoute, useRouter, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { backend, type Profile } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Field, FieldDescription, FieldLabel } from "@/components/ui/field";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

export const Route = createFileRoute("/$locale/$slug/settings/preferences")({
  component: PreferencesPage,
});

function PreferencesPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const router = useRouter();

  // Get profile from parent settings route loader
  const { profile: initialProfile } = settingsRoute.useLoaderData();

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [currentProfile, setCurrentProfile] = React.useState<Profile>(initialProfile);
  const [hideRelations, setHideRelations] = React.useState(initialProfile.hide_relations ?? false);
  const [hideLinks, setHideLinks] = React.useState(initialProfile.hide_links ?? false);

  // Update local state when loader data changes
  React.useEffect(() => {
    setCurrentProfile(initialProfile);
    setHideRelations(initialProfile.hide_relations ?? false);
    setHideLinks(initialProfile.hide_links ?? false);
  }, [initialProfile]);

  // Check if there are unsaved changes
  const hasChanges = React.useMemo(() => {
    return (
      hideRelations !== (initialProfile.hide_relations ?? false) ||
      hideLinks !== (initialProfile.hide_links ?? false)
    );
  }, [hideRelations, hideLinks, initialProfile]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const result = await backend.updateProfile(params.locale, params.slug, {
        hide_relations: hideRelations,
        hide_links: hideLinks,
      });

      if (result === null) {
        toast.error(t("Profile.Failed to update profile"));
        return;
      }

      toast.success(t("Profile.Preferences saved"));

      // Update local state with new values
      setCurrentProfile((prev) => ({
        ...prev,
        hide_relations: hideRelations,
        hide_links: hideLinks,
      }));

      // Invalidate router cache to refresh profile data
      router.invalidate();
    } catch {
      toast.error(t("Profile.Failed to update profile"));
    } finally {
      setIsSubmitting(false);
    }
  };

  // Determine which preferences to show based on profile kind
  const showRelationsPreference = currentProfile.kind === "organization" || currentProfile.kind === "product";
  const showLinksPreference = currentProfile.kind === "organization" || currentProfile.kind === "product";

  // For individual profiles, show contributions preference
  const showContributionsPreference = currentProfile.kind === "individual";

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-serif text-xl font-semibold text-foreground">
            {t("Profile.Preferences")}
          </h3>
          <p className="mt-1 text-sm text-muted-foreground">
            {t("Profile.Control what appears on your profile sidebar.")}
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="mt-6 space-y-6">
        {/* Hide Relations (Members for org/product, Contributions for individual) */}
        {(showRelationsPreference || showContributionsPreference) && (
          <Field orientation="horizontal" className="flex items-center justify-between">
            <div className="space-y-0.5">
              <FieldLabel>
                {showContributionsPreference
                  ? t("Profile.Hide Contributions")
                  : t("Profile.Hide Members")}
              </FieldLabel>
              <FieldDescription>
                {showContributionsPreference
                  ? t("Profile.Hide the Contributions section from your profile sidebar.")
                  : t("Profile.Hide the Members section from your profile sidebar.")}
              </FieldDescription>
            </div>
            <Switch
              checked={hideRelations}
              onCheckedChange={setHideRelations}
            />
          </Field>
        )}

        {/* Hide Links */}
        {showLinksPreference && (
          <Field orientation="horizontal" className="flex items-center justify-between">
            <div className="space-y-0.5">
              <FieldLabel>{t("Profile.Hide Links")}</FieldLabel>
              <FieldDescription>
                {t("Profile.Hide the Links section from your profile sidebar.")}
              </FieldDescription>
            </div>
            <Switch
              checked={hideLinks}
              onCheckedChange={setHideLinks}
            />
          </Field>
        )}

        {/* Save Button */}
        <div className="flex justify-end pt-4">
          <Button type="submit" disabled={!hasChanges || isSubmitting}>
            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isSubmitting ? t("Common.Saving...") : t("Profile.Save Changes")}
          </Button>
        </div>
      </form>
    </Card>
  );
}
