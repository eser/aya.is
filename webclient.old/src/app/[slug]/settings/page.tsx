"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { backend } from "@/shared/modules/backend/backend.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card.tsx";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form.tsx";
import { Input } from "@/shared/components/ui/input.tsx";
import { Textarea } from "@/shared/components/ui/textarea.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Alert, AlertDescription } from "@/shared/components/ui/alert.tsx";
import { ProfilePictureUpload } from "./_components/profile-picture-upload.tsx";
import type { GetProfileData } from "@/shared/modules/backend/profiles/get-profile.ts";
import styles from "./page.module.css";

const generalSettingsSchema = z.object({
  title: z.string().min(2, "Title must be at least 2 characters").max(100, "Title must be no more than 100 characters"),
  description: z.string().min(10, "Description must be at least 10 characters").max(
    500,
    "Description must be no more than 500 characters",
  ),
  pronouns: z.string().max(50, "Pronouns must be no more than 50 characters").optional(),
});

type GeneralSettingsFormData = z.infer<typeof generalSettingsSchema>;

export default function ProfileSettingsPage() {
  const params = useParams();
  const { t, localeCode } = useTranslations();
  const [profile, setProfile] = useState<GetProfileData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const slug = params.slug as string;

  const form = useForm<GeneralSettingsFormData>({
    resolver: zodResolver(generalSettingsSchema),
    defaultValues: {
      title: "",
      description: "",
      pronouns: "",
    },
  });

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const profileData = await backend.getProfile(localeCode, slug);
        if (profileData) {
          setProfile(profileData);
          // Set form values with current profile data
          form.reset({
            title: profileData.title,
            description: profileData.description,
            pronouns: profileData.pronouns || "",
          });
        }
      } catch (error) {
        console.error("Failed to load profile:", error);
        setError(t("Profile", "Failed to load profile data") ?? "Failed to load profile data");
      } finally {
        setIsLoading(false);
      }
    };

    if (slug) {
      loadProfile();
    }
  }, [slug, localeCode, form, t]);

  const onSubmit = async (data: GeneralSettingsFormData) => {
    setIsSubmitting(true);
    setError(null);
    setSuccessMessage(null);

    try {
      // Update profile main fields (excluding title/description which need separate endpoint)
      await backend.updateProfile(localeCode, slug, {
        pronouns: data.pronouns || null,
        // Note: profile_picture_uri is handled by the upload component
        properties: profile?.properties || null,
      });

      // Update translation for current locale
      await backend.updateProfileTranslation(localeCode, slug, localeCode, {
        title: data.title,
        description: data.description,
        properties: null,
      });

      setSuccessMessage(t("Profile", "Profile updated successfully") ?? "Profile updated successfully");

      // Refresh profile data
      const updatedProfile = await backend.getProfile(localeCode, slug);
      if (updatedProfile) {
        setProfile(updatedProfile);
      }
    } catch (error) {
      console.error("Failed to update profile:", error);
      setError(t("Profile", "Failed to update profile") ?? "Failed to update profile");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleProfilePictureUpdate = async (newProfilePictureURI: string) => {
    try {
      await backend.updateProfile(localeCode, slug, {
        profile_picture_uri: newProfilePictureURI === "" ? null : newProfilePictureURI,
        pronouns: profile?.pronouns || null,
        properties: profile?.properties || null,
      });

      // Update local profile state
      if (profile) {
        setProfile({ ...profile, profile_picture_uri: newProfilePictureURI === "" ? null : newProfilePictureURI });
      }

      if (newProfilePictureURI === "") {
        setSuccessMessage(
          t("Profile", "Profile picture removed successfully") ?? "Profile picture removed successfully",
        );
      } else {
        setSuccessMessage(
          t("Profile", "Profile picture updated successfully") ?? "Profile picture updated successfully",
        );
      }
    } catch (error) {
      console.error("Failed to update profile picture:", error);
      setError(t("Profile", "Failed to update profile picture") ?? "Failed to update profile picture");
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loading}>
        {t("Loading", "Loading...") ?? "Loading..."}
      </div>
    );
  }

  if (!profile) {
    return (
      <div className={styles.error}>
        <Alert variant="destructive">
          <AlertDescription>
            {t("Profile", "Profile not found") ?? "Profile not found"}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>
          {t("Profile", "General Settings") ?? "General Settings"}
        </h1>
        <p className={styles.description}>
          {t("Profile", "Update your profile information and preferences.") ??
            "Update your profile information and preferences."}
        </p>
      </div>

      {error && (
        <Alert variant="destructive" className={styles.alert}>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {successMessage && (
        <Alert className={styles.alert}>
          <AlertDescription>{successMessage}</AlertDescription>
        </Alert>
      )}

      <div className={styles.forms}>
        {/* Profile Picture Section */}
        <Card>
          <CardHeader>
            <CardTitle>{t("Profile", "Profile Picture") ?? "Profile Picture"}</CardTitle>
            <CardDescription>
              {t("Profile", "Upload a profile picture to personalize your profile.") ??
                "Upload a profile picture to personalize your profile."}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ProfilePictureUpload
              currentProfilePictureURI={profile.profile_picture_uri}
              onUploadComplete={handleProfilePictureUpdate}
              slug={slug}
              locale={localeCode}
            />
          </CardContent>
        </Card>

        {/* Basic Information Section */}
        <Card>
          <CardHeader>
            <CardTitle>{t("Profile", "Basic Information") ?? "Basic Information"}</CardTitle>
            <CardDescription>
              {t("Profile", "Update your profile title, description and other basic details.") ||
                "Update your profile title, description and other basic details."}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className={styles.form}>
                <FormField
                  control={form.control}
                  name="title"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("Profile", "Title") || "Title"} *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("Profile", "Your name or organization title") ||
                            "Your name or organization title"}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t("Profile", "This is the display name for your profile.") ||
                          "This is the display name for your profile."}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("Profile", "Description") || "Description"} *</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t(
                            "Profile",
                            "Tell people about yourself, your organization, or your product...",
                          ) || "Tell people about yourself, your organization, or your product..."}
                          rows={4}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t("Profile", "A brief description that will appear on your profile page.") ||
                          "A brief description that will appear on your profile page."}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="pronouns"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("Profile", "Pronouns") || "Pronouns"}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("Profile", "e.g., they/them, she/her, he/him") ||
                            "e.g., they/them, she/her, he/him"}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t("Profile", "Optional pronouns to help others address you correctly.") ||
                          "Optional pronouns to help others address you correctly."}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className={styles.actions}>
                  <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting
                      ? (t("Profile", "Saving...") || "Saving...")
                      : (t("Profile", "Save Changes") || "Save Changes")}
                  </Button>
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>

        {/* Read-only Information Section */}
        <Card>
          <CardHeader>
            <CardTitle>{t("Profile", "Profile Information") || "Profile Information"}</CardTitle>
            <CardDescription>
              {t("Profile", "Information about your profile that cannot be changed.") ||
                "Information about your profile that cannot be changed."}
            </CardDescription>
          </CardHeader>
          <CardContent className={styles.readOnlySection}>
            <div className={styles.readOnlyField}>
              <label className={styles.readOnlyLabel}>
                {t("Profile", "Profile URL") || "Profile URL"}
              </label>
              <div className={styles.readOnlyValue}>
                aya.is/{profile.slug}
              </div>
            </div>

            <div className={styles.readOnlyField}>
              <label className={styles.readOnlyLabel}>
                {t("Profile", "Profile Type") || "Profile Type"}
              </label>
              <div className={styles.readOnlyValue}>
                {profile.kind}
              </div>
            </div>

            {profile.custom_domain && (
              <div className={styles.readOnlyField}>
                <label className={styles.readOnlyLabel}>
                  {t("Profile", "Custom Domain") || "Custom Domain"}
                </label>
                <div className={styles.readOnlyValue}>
                  {profile.custom_domain}
                </div>
              </div>
            )}

            <div className={styles.readOnlyField}>
              <label className={styles.readOnlyLabel}>
                {t("Profile", "Created") || "Created"}
              </label>
              <div className={styles.readOnlyValue}>
                {new Date(profile.created_at).toLocaleDateString()}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
