// Profile settings preferences
import * as React from "react";
import { createFileRoute, useRouter, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { backend, type Profile } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Field, FieldLabel } from "@/components/ui/field";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

export const Route = createFileRoute("/$locale/$slug/settings/preferences")({
  component: PreferencesPage,
});

type ModuleVisibility = "public" | "hidden" | "disabled";

function VisibilityRadioGroup(props: {
  value: ModuleVisibility;
  onChange: (value: ModuleVisibility) => void;
  t: (key: string) => string;
  name: string;
}) {
  return (
    <RadioGroup
      value={props.value}
      onValueChange={(val: ModuleVisibility) => props.onChange(val)}
      className="gap-3"
    >
      <div className="flex items-start gap-2">
        <RadioGroupItem value="public" id={`${props.name}-public`} className="mt-0.5" />
        <div>
          <Label htmlFor={`${props.name}-public`} className="font-medium cursor-pointer">
            {props.t("Profile.Public")}
          </Label>
          <p className="text-xs text-muted-foreground">
            {props.t("Profile.Visible in navigation and accessible by everyone.")}
          </p>
        </div>
      </div>
      <div className="flex items-start gap-2">
        <RadioGroupItem value="hidden" id={`${props.name}-hidden`} className="mt-0.5" />
        <div>
          <Label htmlFor={`${props.name}-hidden`} className="font-medium cursor-pointer">
            {props.t("Profile.Hidden")}
          </Label>
          <p className="text-xs text-muted-foreground">
            {props.t("Profile.Not shown in navigation, but accessible via direct link.")}
          </p>
        </div>
      </div>
      <div className="flex items-start gap-2">
        <RadioGroupItem value="disabled" id={`${props.name}-disabled`} className="mt-0.5" />
        <div>
          <Label htmlFor={`${props.name}-disabled`} className="font-medium cursor-pointer">
            {props.t("Profile.Disabled")}
          </Label>
          <p className="text-xs text-muted-foreground">
            {props.t("Profile.Page is completely disabled and returns 404.")}
          </p>
        </div>
      </div>
    </RadioGroup>
  );
}

function PreferencesPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const router = useRouter();

  // Get profile from parent settings route loader
  const { profile: initialProfile } = settingsRoute.useLoaderData();

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [currentProfile, setCurrentProfile] = React.useState<Profile>(initialProfile);
  const [featureRelations, setFeatureRelations] = React.useState<ModuleVisibility>(
    (initialProfile.feature_relations as ModuleVisibility) ?? "public",
  );
  const [featureLinks, setFeatureLinks] = React.useState<ModuleVisibility>(
    (initialProfile.feature_links as ModuleVisibility) ?? "public",
  );
  const [featureQA, setFeatureQA] = React.useState<ModuleVisibility>(
    (initialProfile.feature_qa as ModuleVisibility) ?? "public",
  );

  // Update local state when loader data changes
  React.useEffect(() => {
    setCurrentProfile(initialProfile);
    setFeatureRelations((initialProfile.feature_relations as ModuleVisibility) ?? "public");
    setFeatureLinks((initialProfile.feature_links as ModuleVisibility) ?? "public");
    setFeatureQA((initialProfile.feature_qa as ModuleVisibility) ?? "public");
  }, [initialProfile]);

  // Check if there are unsaved changes
  const hasChanges = React.useMemo(() => {
    return (
      featureRelations !== ((initialProfile.feature_relations as ModuleVisibility) ?? "public") ||
      featureLinks !== ((initialProfile.feature_links as ModuleVisibility) ?? "public") ||
      featureQA !== ((initialProfile.feature_qa as ModuleVisibility) ?? "public")
    );
  }, [featureRelations, featureLinks, featureQA, initialProfile]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const result = await backend.updateProfile(params.locale, params.slug, {
        feature_relations: featureRelations,
        feature_links: featureLinks,
        feature_qa: featureQA,
      });

      if (result === null) {
        toast.error(t("Profile.Failed to update profile"));
        return;
      }

      toast.success(t("Profile.Preferences saved"));

      // Update local state with new values
      setCurrentProfile((prev) => ({
        ...prev,
        feature_relations: featureRelations,
        feature_links: featureLinks,
        feature_qa: featureQA,
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
            {t("Profile.Control how modules appear on your profile.")}
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="mt-6 space-y-8">
        {/* Relations (Members for org/product, Contributions for individual) */}
        {(showRelationsPreference || showContributionsPreference) && (
          <Field>
            <FieldLabel className="text-base font-medium">
              {showContributionsPreference
                ? t("Profile.Contributions Visibility")
                : t("Profile.Members Visibility")}
            </FieldLabel>
            <div className="mt-2">
              <VisibilityRadioGroup
                value={featureRelations}
                onChange={setFeatureRelations}
                t={t}
                name="relations"
              />
            </div>
          </Field>
        )}

        {/* Links */}
        {showLinksPreference && (
          <Field>
            <FieldLabel className="text-base font-medium">
              {t("Profile.Links Visibility")}
            </FieldLabel>
            <div className="mt-2">
              <VisibilityRadioGroup
                value={featureLinks}
                onChange={setFeatureLinks}
                t={t}
                name="links"
              />
            </div>
          </Field>
        )}

        {/* Q&A */}
        <Field>
          <FieldLabel className="text-base font-medium">
            {t("Profile.Q&A Visibility")}
          </FieldLabel>
          <div className="mt-2">
            <VisibilityRadioGroup
              value={featureQA}
              onChange={setFeatureQA}
              t={t}
              name="qa"
            />
          </div>
        </Field>

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
