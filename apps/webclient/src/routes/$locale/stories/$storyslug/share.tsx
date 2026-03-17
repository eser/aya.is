// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Share wizard page for stories
import * as React from "react";
import { createFileRoute, notFound, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend.ts";
import { ShareWizard } from "@/components/share-wizard/index.ts";
import { useAuth } from "@/lib/auth/auth-context.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { Button } from "@/components/ui/button.tsx";
import { PageLayout } from "@/components/page-layouts/default";
import { ShieldX } from "lucide-react";
import { PageNotFound } from "@/components/page-not-found";
import { siteConfig } from "@/config";

export const Route = createFileRoute("/$locale/stories/$storyslug/share")({
  ssr: false,
  loader: async ({ params }) => {
    const { locale, storyslug } = params;

    const story = await backend.getStory(locale, storyslug);
    if (story === null) {
      throw notFound();
    }

    const currentUrl = `${siteConfig.host}/${locale}/stories/${storyslug}`;

    return { story, currentUrl };
  },
  component: ShareWizardPage,
  notFoundComponent: PageNotFound,
});

function ShareWizardPage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { t } = useTranslation();
  const { story, currentUrl } = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState<boolean | null>(null);

  const authorProfileSlug = story.author_profile?.slug ?? null;

  // Check permissions client-side
  React.useEffect(() => {
    if (auth.isLoading) return;

    if (!auth.isAuthenticated) {
      setCanEdit(false);
      return;
    }

    if (authorProfileSlug === null) {
      setCanEdit(false);
      return;
    }

    backend.getStoryPermissions(params.locale, authorProfileSlug, story.id).then((perms) => {
      setCanEdit(perms?.can_edit ?? false);
    });
  }, [auth.isAuthenticated, auth.isLoading, params.locale, authorProfileSlug, story.id]);

  const handleBack = () => {
    navigate({
      to: "/$locale/stories/$storyslug",
      params: {
        locale: params.locale,
        storyslug: params.storyslug,
      },
    });
  };

  // Still checking permissions — show skeleton
  if (canEdit === null) {
    return (
      <PageLayout>
        <div className="flex h-[calc(100vh-140px)] flex-col">
          <div className="flex items-center justify-between border-b p-4">
            <div className="flex items-center gap-3">
              <Skeleton className="h-9 w-20" />
              <Skeleton className="h-6 w-48" />
            </div>
            <Skeleton className="h-9 w-28" />
          </div>
          <div className="flex flex-1 overflow-hidden">
            <div className="flex-1 p-6 space-y-4">
              <Skeleton className="h-6 w-24" />
              <Skeleton className="h-44 w-full" />
              <Skeleton className="h-8 w-64" />
            </div>
            <div className="w-[420px] border-l p-4 space-y-4">
              <Skeleton className="h-6 w-20" />
              <Skeleton className="h-32 w-full" />
            </div>
          </div>
        </div>
      </PageLayout>
    );
  }

  // No permission — show access denied
  if (canEdit !== true || authorProfileSlug === null) {
    return (
      <PageLayout>
        <div className="flex flex-col items-center justify-center gap-4 py-16 text-center">
          <ShieldX className="size-16 text-muted-foreground" />
          <h2 className="text-xl font-semibold">
            {t("ShareWizard.Access Denied")}
          </h2>
          <p className="text-muted-foreground max-w-md">
            {t("ShareWizard.You don't have permission to use the share wizard for this story.")}
          </p>
          <Button variant="outline" onClick={handleBack}>
            {t("ShareWizard.Back to Story")}
          </Button>
        </div>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <ShareWizard
        story={story}
        locale={params.locale}
        currentUrl={currentUrl}
        onBack={handleBack}
      />
    </PageLayout>
  );
}
