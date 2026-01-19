// Profile settings index (general settings)
import * as React from "react";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { backend, type Profile } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { ProfilePictureUpload } from "@/components/profile/profile-picture-upload";
import { EditProfileForm, type EditProfileInput } from "@/components/forms/edit-profile-form";

// Import parent route to access its loader data
import { Route as SettingsRoute } from "./route";

export const Route = createFileRoute("/$locale/$slug/settings/")({
  component: SettingsIndexPage,
});

function SettingsIndexPage() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const router = useRouter();

  // Get profile from parent settings route loader
  const { profile: initialProfile } = SettingsRoute.useLoaderData();

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [currentProfile, setCurrentProfile] = React.useState<Profile>(initialProfile);

  // Update local profile state when loader data changes
  React.useEffect(() => {
    setCurrentProfile(initialProfile);
  }, [initialProfile]);

  const handleSubmit = async (data: EditProfileInput) => {
    setIsSubmitting(true);
    try {
      // Update pronouns via profile endpoint (only if individual profile)
      if (currentProfile.kind === "individual") {
        const profileResult = await backend.updateProfile(params.locale, params.slug, {
          pronouns: data.pronouns || null,
        });

        if (profileResult === null) {
          toast.error(t("Profile.Failed to update profile"));
          return;
        }
      }

      // Update title/description via translation endpoint
      const translationResult = await backend.updateProfileTranslation(
        params.locale,
        params.slug,
        params.locale, // Update current locale's translation
        {
          title: data.title,
          description: data.description,
        },
      );

      if (translationResult === null) {
        toast.error(t("Profile.Failed to update profile"));
        return;
      }

      toast.success(t("Profile.Profile updated successfully"));

      // Update local state with new values
      setCurrentProfile((prev) => ({
        ...prev,
        title: data.title,
        description: data.description,
        pronouns: data.pronouns || null,
      }));

      // Invalidate router cache to refresh profile data
      router.invalidate();
    } catch {
      toast.error(t("Profile.Failed to update profile"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePictureUpload = (newUri: string) => {
    // Update local state with new picture URI
    setCurrentProfile((prev) => ({
      ...prev,
      profile_picture_uri: newUri,
    }));

    // Invalidate router cache to refresh profile data
    router.invalidate();
  };

  return (
    <Card className="p-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-serif text-xl font-semibold text-foreground">{t("Profile.General")}</h3>
          {/* <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Manage your social media links and external websites.")}
          </p> */}
        </div>
        {/* <Button onClick={handleOpenAddDialog}>
          <Plus className="size-4 mr-2" />
          {t("Profile.Add Link")}
        </Button> */}
      </div>

      <div className="space-y-8">
        {/* Profile Picture Section */}
        <div>
          <ProfilePictureUpload
            currentImageUri={currentProfile.profile_picture_uri}
            profileSlug={params.slug}
            profileTitle={currentProfile.title}
            locale={params.locale}
            onUploadComplete={handlePictureUpload}
          />
        </div>

        {/* Divider */}
        <hr className="border-border" />

        {/* Profile Form Section */}
        <div>
          <EditProfileForm
            profile={currentProfile}
            onSubmit={handleSubmit}
            isSubmitting={isSubmitting}
          />
        </div>
      </div>
    </Card>
  );
}
