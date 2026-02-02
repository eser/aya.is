import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import { useAuth } from "@/lib/auth/auth-context";
import { PageLayout } from "@/components/page-layouts/default";
import { CreateProfileForm } from "@/components/forms/create-profile-form";
import type { CreateProfileInput } from "@/lib/schemas/profile";

export const Route = createFileRoute("/$locale/products/new")({
  component: NewProductProfilePage,
});

function NewProductProfilePage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { t } = useTranslation();

  if (auth.isLoading) {
    return (
      <PageLayout>
        <div className="flex items-center justify-center py-20">
          <div className="text-muted-foreground">{t("Loading.Loading...")}</div>
        </div>
      </PageLayout>
    );
  }

  if (!auth.isAuthenticated || auth.user === null) {
    return (
      <PageLayout>
        <div className="content py-20 text-center">
          <h2>{t("Auth.Access Denied")}</h2>
          <p>{t("Auth.You need to be logged in to create a profile.")}</p>
        </div>
      </PageLayout>
    );
  }

  const hasIndividualProfile = auth.user.individual_profile_slug !== undefined;

  const handleSubmit = async (data: CreateProfileInput) => {
    try {
      const result = await backend.createProfile(params.locale, {
        kind: data.kind,
        slug: data.slug,
        title: data.title,
        description: data.description ?? "",
      });

      if (result !== null) {
        toast.success(t("Profile.Profile created successfully"));
        navigate({
          to: "/$locale/$slug",
          params: { locale: params.locale, slug: data.slug },
        });
      } else {
        toast.error(t("Profile.Failed to create profile"));
      }
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error(t("Profile.Failed to create profile"));
      }
    }
  };

  return (
    <PageLayout>
      <CreateProfileForm
        locale={params.locale}
        defaultKind="product"
        backUrl={`/${params.locale}/products`}
        hasIndividualProfile={hasIndividualProfile}
        onSubmit={handleSubmit}
      />
    </PageLayout>
  );
}
