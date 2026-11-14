// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Privacy settings: download personal data, hide the profile, delete the account
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/lib/auth/auth-context";
import { privacyStateQueryOptions } from "@/modules/backend/queries";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Skeleton } from "@/components/ui/skeleton";
import { DeleteSection } from "@/components/privacy/delete-section";
import { ExportSection } from "@/components/privacy/export-section";
import { HideSection } from "@/components/privacy/hide-section";

export const Route = createFileRoute("/$locale/$slug/settings/privacy")({
  component: PrivacySettingsPage,
});

const settingsRoute = getRouteApi("/$locale/$slug/settings");

function PrivacySettingsPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const queryClient = useQueryClient();
  const { user, isLoading: isAuthLoading } = useAuth();

  // Ensure parent loader data is available (provides profile context)
  settingsRoute.useLoaderData();

  // These actions operate on the signed-in user's own account, so they only make sense on their own profile.
  const isOwnAccount = user !== null && user.individual_profile_slug === params.slug;

  const { data: state, isLoading: isStateLoading } = useQuery({
    ...privacyStateQueryOptions(params.locale),
    enabled: isOwnAccount,
  });

  if (isAuthLoading || (isOwnAccount && isStateLoading)) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-48 w-full" />
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!isOwnAccount || state === null || state === undefined) {
    return (
      <Alert>
        <AlertDescription>{t("Privacy.Privacy settings are only available for your own profile.")}</AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <p className="text-muted-foreground">
        {t("Privacy.Manage how your personal data is stored, shown and removed.")}
      </p>

      <ExportSection
        locale={params.locale}
        state={state.export}
        onChanged={() => {
          void queryClient.invalidateQueries({ queryKey: ["privacy", params.locale] });
        }}
      />

      <HideSection locale={params.locale} />

      <DeleteSection locale={params.locale} slug={state.profile_slug} />
    </div>
  );
}
