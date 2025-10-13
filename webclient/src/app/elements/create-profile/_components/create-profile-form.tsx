"use client";

import { useCallback, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/shared/components/ui/select.tsx";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card.tsx";
import { Alert, AlertDescription } from "@/shared/components/ui/alert.tsx";
import { forbiddenSlugs, supportedLocales } from "@/shared/config.ts";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
const forbiddenPattern = `^(?!${forbiddenSlugs.join("|")}).*`;
const forbiddenRegex = new RegExp(forbiddenPattern, "i");
const createProfileSchema = z.object({
  kind: z.enum(["individual", "organization", "product"], {
    message: "Profile type is required",
  }),
  slug: z
    .string()
    .min(3, "Slug must be at least 3 characters long")
    .max(50, "Slug must be no more than 50 characters long")
    .regex(/^[a-z0-9-]+$/, "Slug can only contain lowercase letters, numbers, and hyphens")
    .regex(forbiddenRegex, "This slug is reserved and cannot be used"),
  title: z
    .string()
    .min(2, "Title must be at least 2 characters long")
    .max(100, "Title must be no more than 100 characters long"),
  description: z
    .string()
    .min(10, "Description must be at least 10 characters long")
    .max(500, "Description must be no more than 500 characters long"),
  locale: z.string().min(1, "Locale is required"),
});

type CreateProfileFormData = z.infer<typeof createProfileSchema>;

type SlugAvailabilityState = {
  isChecking: boolean;
  isAvailable: boolean | null;
  error: string | null;
};

export function CreateProfileForm() {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [slugAvailability, setSlugAvailability] = useState<SlugAvailabilityState>({
    isChecking: false,
    isAvailable: null,
    error: null,
  });
  const router = useRouter();
  const { session, isLoading } = useAuth();
  const { localeCode } = useTranslations();

  // Use the proper locale from translations context
  const currentLocale = localeCode;

  // Check if session already has an individual profile
  const hasIndividualProfile = session?.selected_profile?.kind === "individual";
  const defaultKind = hasIndividualProfile ? "organization" : "individual";

  const form = useForm<CreateProfileFormData>({
    resolver: zodResolver(createProfileSchema),
    defaultValues: {
      kind: defaultKind,
      slug: "",
      title: "",
      description: "",
      locale: currentLocale,
    },
  });

  // Update form values when session data or locale changes
  useEffect(() => {
    const currentKind = form.getValues("kind");
    const currentFormLocale = form.getValues("locale");

    // Update kind if needed
    if (currentKind !== defaultKind) {
      form.setValue("kind", defaultKind, { shouldDirty: true });
    }

    // Update locale if needed
    if (currentFormLocale !== currentLocale) {
      form.setValue("locale", currentLocale, { shouldDirty: true });
    }
  }, [session, hasIndividualProfile, currentLocale, form]);

  const checkSlugAvailability = useCallback(async (slug: string) => {
    if (slug === "" || slug.length < 3) {
      setSlugAvailability({ isChecking: false, isAvailable: null, error: null });
      return;
    }

    // Check against forbidden slugs first
    if (forbiddenSlugs.includes(slug.toLowerCase())) {
      setSlugAvailability({
        isChecking: false,
        isAvailable: false,
        error: "This slug is reserved and cannot be used.",
      });
      return;
    }

    setSlugAvailability({ isChecking: true, isAvailable: null, error: null });

    try {
      const exists = await backend.checkProfileSlug(currentLocale, slug);
      const isAvailable = !exists;

      setSlugAvailability({
        isChecking: false,
        isAvailable,
        error: null,
      });
    } catch (_err: unknown) {
      setSlugAvailability({
        isChecking: false,
        isAvailable: null,
        error: "Unable to check slug availability. Please try again.",
      });
    }
  }, []);

  // Debounced slug checking
  useEffect(() => {
    const slug = form.watch("slug");
    const timeoutId = setTimeout(() => {
      checkSlugAvailability(slug);
    }, 500); // 500ms debounce

    return () => clearTimeout(timeoutId);
  }, [form.watch("slug"), checkSlugAvailability]);

  const onSubmit = async (data: CreateProfileFormData) => {
    // Check if user is trying to create individual profile when they already have one
    if (data.kind === "individual" && hasIndividualProfile) {
      setError("You already have an individual profile. Please select a different profile type.");
      return;
    }

    // Check slug availability one final time before submission
    if (slugAvailability.isAvailable === false) {
      setError("Cannot create profile: the selected slug is not available.");
      return;
    }

    if (slugAvailability.isAvailable === null) {
      setError("Please wait for slug availability check to complete.");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      // Call the backend function to create the profile
      const profileData = await backend.createProfile(data.locale, {
        kind: data.kind,
        slug: data.slug,
        title: data.title,
        description: data.description,
      });

      if (profileData === null) {
        throw new Error("Failed to create profile: No data returned");
      }

      // Redirect to the created profile page
      router.push(`/${profileData.slug}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred while creating the profile");
    } finally {
      setIsSubmitting(false);
    }
  };

  const watchedSlug = form.watch("slug", "");

  // Don't render form until user data is loaded
  if (isLoading) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center py-8">
          <div className="text-center">Loading user information...</div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Profile Information</CardTitle>
        <CardDescription>
          Fill in the details for your new profile. Required fields are marked with *.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            {hasIndividualProfile && (
              <Alert>
                <AlertDescription>
                  <strong>Notes:</strong>
                  <ul className="list-disc list-inside">
                    <li>
                      You already have an individual profile. You can create additional organization or product
                      profiles, but each user can only have one individual profile.
                    </li>
                    <li>
                      Profile will be created in the selected locale ({supportedLocales[form.watch("locale")]?.name}).
                      After creation, you can add translations for additional languages through created profile
                      settings.
                    </li>
                  </ul>
                </AlertDescription>
              </Alert>
            )}

            <FormField
              control={form.control}
              name="locale"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Language / Locale *</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    value={field.value}
                    key={`locale-${field.value}`}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select a language" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {Object.values(supportedLocales).map((locale) => (
                        <SelectItem key={locale.code} value={locale.code}>
                          {locale.flag} {locale.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    The profile will be created in this language. You can add additional languages later.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="kind"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Profile Type *</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    value={field.value}
                    key={`profile-type-${field.value}`}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select profile type" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem
                        value="individual"
                        disabled={hasIndividualProfile}
                        className={hasIndividualProfile ? "opacity-50 cursor-not-allowed" : ""}
                      >
                        👤 Individual {hasIndividualProfile ? "(Already exists)" : ""}
                      </SelectItem>
                      <SelectItem value="organization">🏢 Organization</SelectItem>
                      <SelectItem value="product">📦 Product</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    Choose the type of profile you're creating.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="slug"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Profile URL (Slug) *</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        placeholder="your-profile-name"
                        {...field}
                        onChange={(e) => {
                          // Convert to lowercase and replace invalid characters
                          const value = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "-");
                          field.onChange(value);
                        }}
                        className={watchedSlug.length >= 3
                          ? slugAvailability.isAvailable === true
                            ? "border-green-500 focus:border-green-500"
                            : slugAvailability.isAvailable === false
                            ? "border-red-500 focus:border-red-500"
                            : ""
                          : ""}
                      />
                      {watchedSlug.length >= 3 && (
                        <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                          {slugAvailability.isChecking && (
                            <div className="w-4 h-4 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin" />
                          )}
                          {!slugAvailability.isChecking && slugAvailability.isAvailable === true && (
                            <span className="text-green-500 text-sm">✓</span>
                          )}
                          {!slugAvailability.isChecking && slugAvailability.isAvailable === false && (
                            <span className="text-red-500 text-sm">✗</span>
                          )}
                        </div>
                      )}
                    </div>
                  </FormControl>
                  <FormDescription>
                    {watchedSlug.length >= 3 && slugAvailability.isChecking && (
                      <span className="text-blue-600">Checking availability...</span>
                    )}
                    {watchedSlug.length >= 3 && !slugAvailability.isChecking && slugAvailability.isAvailable === true &&
                      <span className="text-green-600">✓ Available! Your profile will be at: aya.is/{watchedSlug}
                      </span>}
                    {watchedSlug.length >= 3 && !slugAvailability.isChecking &&
                      slugAvailability.isAvailable === false && (
                      <span className="text-red-600">✗ {slugAvailability.error ?? "This slug is not available"}</span>
                    )}
                    {watchedSlug.length >= 3 && !slugAvailability.isChecking && slugAvailability.isAvailable === null &&
                      <span className="text-gray-600">Your profile will be available at: aya.is/{watchedSlug}</span>}
                    {watchedSlug.length < 3 && "This will be used in your profile URL (e.g., aya.is/your-profile-name)"}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Title *</FormLabel>
                  <FormControl>
                    <Input placeholder="Your name or organization title" {...field} />
                  </FormControl>
                  <FormDescription>
                    This is the display name for your profile.
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
                  <FormLabel>Description *</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Tell people about yourself, your organization, or your product..."
                      rows={4}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    A brief description that will appear on your profile page.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex gap-4 justify-end">
              <Button
                type="submit"
                disabled={isSubmitting ||
                  slugAvailability.isChecking ||
                  slugAvailability.isAvailable === false ||
                  (watchedSlug.length >= 3 && slugAvailability.isAvailable === null)}
              >
                {isSubmitting ? "Creating Profile..." : "Create Profile"}
              </Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
