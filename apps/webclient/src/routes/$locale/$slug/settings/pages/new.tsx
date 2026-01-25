// Create new profile page
import { createFileRoute, useNavigate, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import {
  ContentEditor,
  type ContentEditorData,
} from "@/components/content-editor";
import { useAuth } from "@/lib/auth/auth-context";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

export const Route = createFileRoute("/$locale/$slug/settings/pages/new")({
  component: NewPagePage,
});

function NewPagePage() {
  const params = Route.useParams();
  const navigate = useNavigate();
  const auth = useAuth();
  const { t } = useTranslation();
  const { profile } = settingsRoute.useLoaderData();

  // No permission - settings route already handles redirect, but just in case
  if (profile === null) {
    return (
      <div className="content">
        <h2>{t("Auth.Access Denied")}</h2>
        <p>{t("Profile.You don't have permission to create pages for this profile.")}</p>
      </div>
    );
  }

  const initialData: ContentEditorData = {
    title: "",
    slug: "",
    summary: "",
    content: "",
    storyPictureUri: null,
    status: "draft",
  };

  const handleSave = async (data: ContentEditorData) => {
    const result = await backend.createProfilePage(
      params.locale,
      params.slug,
      {
        slug: data.slug,
        title: data.title,
        summary: data.summary,
        content: data.content,
        cover_picture_uri: data.storyPictureUri,
        published_at: data.publishedAt,
      },
    );

    if (result !== null) {
      toast.success(t("Profile.Page created successfully"));
      navigate({
        to: "/$locale/$slug/$pageslug",
        params: {
          locale: params.locale,
          slug: params.slug,
          pageslug: data.slug,
        },
      });
    } else {
      toast.error(t("Profile.Failed to create page"));
    }
  };

  return (
    <div className="h-[calc(100vh-140px)]">
      <ContentEditor
        locale={params.locale}
        profileSlug={params.slug}
        contentType="page"
        initialData={initialData}
        backUrl={`/${params.locale}/${params.slug}/settings/pages`}
        userKind={auth.user?.kind}
        onSave={handleSave}
        isNew
      />
    </div>
  );
}
