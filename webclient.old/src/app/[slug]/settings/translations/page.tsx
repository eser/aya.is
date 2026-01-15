"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { backend } from "@/shared/modules/backend/backend.ts";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
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
import { Badge } from "@/shared/components/ui/badge.tsx";
import type { GetProfileData } from "@/shared/modules/backend/profiles/get-profile.ts";
import type { ProfileTranslation } from "@/shared/modules/backend/profiles/get-profile-translations.ts";
import styles from "./page.module.css";

const translationSchema = z.object({
  title: z.string().min(2, "Title must be at least 2 characters").max(100, "Title must be no more than 100 characters"),
  description: z.string().min(10, "Description must be at least 10 characters").max(
    500,
    "Description must be no more than 500 characters",
  ),
});

type TranslationFormData = z.infer<typeof translationSchema>;

// All supported locales with their display names
const SUPPORTED_LOCALES = [
  { code: "ar", name: "العربية", nativeName: "Arabic" },
  { code: "de", name: "Deutsch", nativeName: "German" },
  { code: "en", name: "English", nativeName: "English" },
  { code: "es", name: "Español", nativeName: "Spanish" },
  { code: "fr", name: "Français", nativeName: "French" },
  { code: "it", name: "Italiano", nativeName: "Italian" },
  { code: "ja", name: "日本語", nativeName: "Japanese" },
  { code: "ko", name: "한국어", nativeName: "Korean" },
  { code: "nl", name: "Nederlands", nativeName: "Dutch" },
  { code: "pt-PT", name: "Português", nativeName: "Portuguese" },
  { code: "ru", name: "Русский", nativeName: "Russian" },
  { code: "tr", name: "Türkçe", nativeName: "Turkish" },
  { code: "zh-CN", name: "中文", nativeName: "Chinese (Simplified)" },
];

export default function ProfileTranslationsPage() {
  const params = useParams();
  const { t, localeCode } = useTranslations();
  const { isAuthenticated } = useAuth();
  const [profile, setProfile] = useState<GetProfileData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [selectedLocale, setSelectedLocale] = useState<string>(localeCode);
  const [allTranslations, setAllTranslations] = useState<Map<string, ProfileTranslation>>(new Map());

  const slug = params.slug as string;

  const form = useForm<TranslationFormData>({
    resolver: zodResolver(translationSchema),
    defaultValues: {
      title: "",
      description: "",
    },
  });

  useEffect(() => {
    const loadProfileAndTranslations = async () => {
      try {
        // Load the main profile data
        const profileData = await backend.getProfile(localeCode, slug);
        if (profileData) {
          setProfile(profileData);

          // Only load translations if user is authenticated
          if (isAuthenticated) {
            // Load all translations at once using the new endpoint
            const translationsResponse = await backend.getProfileTranslations(localeCode, slug);
            if (translationsResponse) {
              const translationsMap = new Map<string, ProfileTranslation>();
              translationsResponse.forEach((translation: ProfileTranslation) => {
                translationsMap.set(translation.locale_code, translation);
              });
              setAllTranslations(translationsMap);
            }
          }
        }
      } catch (error) {
        console.error("Failed to load profile data:", error);
        setError(t("Profile", "Failed to load profile data") ?? "Failed to load profile data");
      } finally {
        setIsLoading(false);
      }
    };

    if (slug) {
      loadProfileAndTranslations();
    }
  }, [slug, localeCode, t, isAuthenticated]);

  useEffect(() => {
    if (!profile || !selectedLocale) return;

    // Use cached translation data instead of making API call
    const translation = allTranslations.get(selectedLocale);
    if (translation) {
      form.reset({
        title: translation.title || "",
        description: translation.description || "",
      });
    } else {
      // No translation exists for this locale, reset form
      form.reset({
        title: "",
        description: "",
      });
    }
  }, [selectedLocale, profile, allTranslations, form]);

  const onSubmit = async (data: TranslationFormData) => {
    setIsSubmitting(true);
    setError(null);
    setSuccessMessage(null);

    try {
      await backend.updateProfileTranslation(localeCode, slug, selectedLocale, {
        title: data.title,
        description: data.description,
        properties: null,
      });

      // Update local cache with the new translation
      setAllTranslations((prev) => {
        const newMap = new Map(prev);
        newMap.set(selectedLocale, {
          profile_id: profile?.id ?? "",
          locale_code: selectedLocale,
          title: data.title,
          description: data.description,
          properties: null,
        });
        return newMap;
      });

      const localeName = SUPPORTED_LOCALES.find((l) => l.code === selectedLocale)?.name ?? selectedLocale;
      setSuccessMessage(
        t("Profile", `Translation for ${localeName} updated successfully`) ??
          `Translation for ${localeName} updated successfully`,
      );
    } catch (error) {
      console.error("Failed to update translation:", error);
      setError(t("Profile", "Failed to update translation") ?? "Failed to update translation");
    } finally {
      setIsSubmitting(false);
    }
  };

  const copyFromLocale = (sourceLocale: string) => {
    const sourceTranslation = allTranslations.get(sourceLocale);
    if (sourceTranslation) {
      form.reset({
        title: sourceTranslation.title || "",
        description: sourceTranslation.description || "",
      });

      const sourceName = SUPPORTED_LOCALES.find((l) => l.code === sourceLocale)?.name ?? sourceLocale;
      setSuccessMessage(
        t("Profile", `Content copied from ${sourceName}`) ??
          `Content copied from ${sourceName}`,
      );
    } else {
      setError(t("Profile", "Failed to copy content") || "Failed to copy content");
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loading}>
        {t("Loading", "Loading...") || "Loading..."}
      </div>
    );
  }

  if (!profile) {
    return (
      <div className={styles.error}>
        <Alert variant="destructive">
          <AlertDescription>
            {t("Profile", "Profile not found") || "Profile not found"}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  const currentLocaleInfo = SUPPORTED_LOCALES.find((l) => l.code === selectedLocale);
  const availableSourceLocales = SUPPORTED_LOCALES.filter((l) =>
    l.code !== selectedLocale && allTranslations.has(l.code)
  );

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>
          {t("Profile", "Translations") || "Translations"}
        </h1>
        <p className={styles.description}>
          {t("Profile", "Manage translations for your profile in different languages.") ||
            "Manage translations for your profile in different languages."}
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

      {/* Language Overview */}
      <Card>
        <CardHeader>
          <CardTitle>{t("Profile", "Available Languages") || "Available Languages"}</CardTitle>
          <CardDescription>
            {t("Profile", "Languages with existing translations are highlighted.") ||
              "Languages with existing translations are highlighted."}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className={styles.localeGrid}>
            {SUPPORTED_LOCALES.map((locale) => (
              <div
                key={locale.code}
                className={`${styles.localeCard} ${allTranslations.has(locale.code) ? styles.localeCardActive : ""}`}
                onClick={() => setSelectedLocale(locale.code)}
              >
                <div className={styles.localeInfo}>
                  <span className={styles.localeCode}>{locale.code}</span>
                  <span className={styles.localeName}>{locale.name}</span>
                </div>
                {allTranslations.has(locale.code) && (
                  <Badge variant="secondary" className={styles.badge}>
                    {t("Profile", "Translated") || "Translated"}
                  </Badge>
                )}
                {locale.code === selectedLocale && (
                  <Badge className={styles.badge}>
                    {t("Profile", "Selected") || "Selected"}
                  </Badge>
                )}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Translation Form */}
      <Card>
        <CardHeader>
          <CardTitle>
            {t("Profile", "Edit Translation") || "Edit Translation"} - {currentLocaleInfo?.name || selectedLocale}
          </CardTitle>
          <CardDescription>
            {t("Profile", "Update the title and description for this language.") ||
              "Update the title and description for this language."}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {availableSourceLocales.length > 0 && (
            <div className={styles.copySection}>
              <label className={styles.copyLabel}>
                {t("Profile", "Copy from existing translation:") || "Copy from existing translation:"}
              </label>
              <div className={styles.copyButtons}>
                {availableSourceLocales.map((locale) => (
                  <Button
                    key={locale.code}
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => copyFromLocale(locale.code)}
                    className={styles.copyButton}
                  >
                    {locale.name} ({locale.code})
                  </Button>
                ))}
              </div>
            </div>
          )}

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
                        placeholder={t("Profile", "Profile title in this language") || "Profile title in this language"}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("Profile", "The display name for your profile in this language.") ||
                        "The display name for your profile in this language."}
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
                        placeholder={t("Profile", "Profile description in this language...") ||
                          "Profile description in this language..."}
                        rows={4}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("Profile", "A brief description that will appear on your profile page in this language.") ||
                        "A brief description that will appear on your profile page in this language."}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className={styles.actions}>
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting
                    ? (t("Profile", "Saving...") || "Saving...")
                    : (t("Profile", "Save Translation") || "Save Translation")}
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
