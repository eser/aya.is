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
import { Button } from "@/shared/components/ui/button.tsx";
import { Alert, AlertDescription } from "@/shared/components/ui/alert.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/shared/components/ui/select.tsx";
import { Checkbox } from "@/shared/components/ui/checkbox.tsx";
import type { ProfileLink } from "@/shared/modules/backend/profiles/list-profile-links.ts";
import styles from "./page.module.css";

const linkSchema = z.object({
  kind: z.string().min(1, "Link type is required"),
  title: z.string().min(2, "Title must be at least 2 characters").max(100, "Title must be no more than 100 characters"),
  uri: z.string().refine((val) => {
    if (!val || val === "") return true;
    try {
      new URL(val);
      return true;
    } catch {
      return false;
    }
  }, "Please enter a valid URL").optional(),
  is_hidden: z.boolean(),
});

type LinkFormData = z.infer<typeof linkSchema>;

// Available link types with their display names and example URLs
const LINK_TYPES = [
  {
    value: "website",
    label: "Website",
    icon: "🌐",
    placeholder: "https://example.com",
  },
  {
    value: "github",
    label: "GitHub",
    icon: "📂",
    placeholder: "https://github.com/username",
  },
  {
    value: "x",
    label: "X (Twitter)",
    icon: "𝕏",
    placeholder: "https://x.com/username",
  },
  {
    value: "bsky",
    label: "Bluesky",
    icon: "☁️",
    placeholder: "https://bsky.app/profile/username",
  },
];

export default function ProfileLinksPage() {
  const params = useParams();
  const { t, localeCode } = useTranslations();
  const [links, setLinks] = useState<ProfileLink[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [editingLink, setEditingLink] = useState<ProfileLink | null>(null);
  const [showForm, setShowForm] = useState(false);

  const slug = params.slug as string;

  const form = useForm<LinkFormData>({
    resolver: zodResolver(linkSchema),
    defaultValues: {
      kind: "",
      title: "",
      uri: "",
      is_hidden: false,
    },
  });

  useEffect(() => {
    const loadLinks = async () => {
      try {
        const linksData = await backend.listProfileLinks(localeCode, slug);
        if (linksData) {
          setLinks(linksData);
        }
      } catch (error) {
        console.error("Failed to load profile links:", error);
        setError(t("Profile", "Failed to load profile links") ?? "Failed to load profile links");
      } finally {
        setIsLoading(false);
      }
    };

    if (slug) {
      loadLinks();
    }
  }, [slug, localeCode, t]);

  const onSubmit = async (data: LinkFormData) => {
    setIsSubmitting(true);
    setError(null);
    setSuccessMessage(null);

    try {
      if (editingLink) {
        // Update existing link
        const updatedLink = await backend.updateProfileLink(localeCode, slug, editingLink.id, {
          ...data,
          uri: data.uri ?? null,
          order: editingLink.order, // Keep existing order for now
        });

        if (updatedLink === null) {
          throw new Error("Failed to update link");
        }

        // Update in local state
        setLinks((prev) => prev.map((link) => link.id === editingLink.id ? updatedLink : link));

        setSuccessMessage(t("Profile", "Link updated successfully") ?? "Link updated successfully");
      } else {
        // Create new link
        const newLink = await backend.createProfileLink(localeCode, slug, {
          ...data,
          uri: data.uri ?? null,
        });

        if (newLink === null) {
          throw new Error("Failed to create link");
        }

        // Add to local state
        setLinks((prev) => [...prev, newLink]);

        setSuccessMessage(t("Profile", "Link created successfully") ?? "Link created successfully");
      }

      // Reset form and close
      form.reset();
      setEditingLink(null);
      setShowForm(false);
    } catch (error) {
      console.error("Failed to save link:", error);
      setError(t("Profile", "Failed to save link") ?? "Failed to save link");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEdit = (link: ProfileLink) => {
    setEditingLink(link);
    form.reset({
      kind: link.kind,
      title: link.title,
      uri: link.uri ?? "",
      is_hidden: link.is_hidden,
    });
    setShowForm(true);
  };

  const handleDelete = async (linkId: string) => {
    if (
      !confirm(
        t("Profile", "Are you sure you want to delete this link?") ?? "Are you sure you want to delete this link?",
      )
    ) {
      return;
    }

    try {
      await backend.deleteProfileLink(localeCode, slug, linkId);

      // Remove from local state
      setLinks((prev) => prev.filter((link) => link.id !== linkId));

      setSuccessMessage(t("Profile", "Link deleted successfully") ?? "Link deleted successfully");
    } catch (error) {
      console.error("Failed to delete link:", error);
      setError(t("Profile", "Failed to delete link") ?? "Failed to delete link");
    }
  };

  const handleAdd = () => {
    setEditingLink(null);
    form.reset({
      kind: "",
      title: "",
      uri: "",
      is_hidden: false,
    });
    setShowForm(true);
  };

  const handleCancel = () => {
    setEditingLink(null);
    setShowForm(false);
    form.reset();
  };

  if (isLoading) {
    return (
      <div className={styles.loading}>
        {t("Loading", "Loading...") ?? "Loading..."}
      </div>
    );
  }

  const selectedLinkType = LINK_TYPES.find((type) => type.value === form.watch("kind"));

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>
          {t("Profile", "Social Links") ?? "Social Links"}
        </h1>
        <p className={styles.description}>
          {t("Profile", "Manage your social media links and external websites.") ??
            "Manage your social media links and external websites."}
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

      {/* Current Links */}
      <Card>
        <CardHeader>
          <div className={styles.cardHeaderWithAction}>
            <div>
              <CardTitle>{t("Profile", "Current Links") ?? "Current Links"}</CardTitle>
              <CardDescription>
                {t("Profile", "Your existing social media links and websites.") ??
                  "Your existing social media links and websites."}
              </CardDescription>
            </div>
            <Button onClick={handleAdd}>
              {t("Profile", "Add Link") ?? "Add Link"}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {links.length === 0
            ? (
              <div className={styles.emptyState}>
                <p className={styles.emptyText}>
                  {t("Profile", "No links added yet.") ?? "No links added yet."}
                </p>
                <Button onClick={handleAdd} variant="outline">
                  {t("Profile", "Add Your First Link") ?? "Add Your First Link"}
                </Button>
              </div>
            )
            : (
              <div className={styles.linksList}>
                {links
                  .sort((a, b) => a.order - b.order)
                  .map((link) => {
                    const linkType = LINK_TYPES.find((type) => type.value === link.kind);
                    return (
                      <div key={link.id} className={styles.linkItem}>
                        <div className={styles.linkInfo}>
                          <div className={styles.linkTypeIcon}>
                            {linkType?.icon ?? "🔗"}
                          </div>
                          <div className={styles.linkDetails}>
                            <div className={styles.linkTitle}>
                              {link.title}
                              {link.is_hidden && (
                                <span className={styles.hiddenBadge}>
                                  {t("Profile", "Hidden") ?? "Hidden"}
                                </span>
                              )}
                            </div>
                            <div className={styles.linkType}>
                              {linkType?.label ?? link.kind}
                            </div>
                            {link.uri && (
                              <div className={styles.linkUrl}>
                                <a
                                  href={link.uri}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className={styles.linkUrlLink}
                                >
                                  {link.uri}
                                </a>
                              </div>
                            )}
                          </div>
                        </div>
                        <div className={styles.linkActions}>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEdit(link)}
                          >
                            {t("Profile", "Edit") ?? "Edit"}
                          </Button>
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDelete(link.id)}
                          >
                            {t("Profile", "Delete") ?? "Delete"}
                          </Button>
                        </div>
                      </div>
                    );
                  })}
              </div>
            )}
        </CardContent>
      </Card>

      {/* Add/Edit Form */}
      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>
              {editingLink
                ? (t("Profile", "Edit Link") ?? "Edit Link")
                : (t("Profile", "Add New Link") ?? "Add New Link")}
            </CardTitle>
            <CardDescription>
              {t("Profile", "Configure your social media link or website.") ??
                "Configure your social media link or website."}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className={styles.form}>
                <FormField
                  control={form.control}
                  name="kind"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("Profile", "Link Type") ?? "Link Type"} *</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder={t("Profile", "Select link type") ?? "Select link type"} />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {LINK_TYPES.map((type) => (
                            <SelectItem key={type.value} value={type.value}>
                              <span className={styles.selectOption}>
                                <span className={styles.selectIcon}>{type.icon}</span>
                                {type.label}
                              </span>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormDescription>
                        {t("Profile", "Choose the type of link you're adding.") ??
                          "Choose the type of link you're adding."}
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
                      <FormLabel>{t("Profile", "Display Title") ?? "Display Title"} *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("Profile", "How this link should appear") ?? "How this link should appear"}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t("Profile", "A friendly name for this link.") ?? "A friendly name for this link."}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="uri"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("Profile", "URL") ?? "URL"}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={selectedLinkType?.placeholder ?? "https://example.com"}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        {t("Profile", "The full URL to your profile or website.") ??
                          "The full URL to your profile or website."}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="is_hidden"
                  render={({ field }) => (
                    <FormItem className={styles.checkboxItem}>
                      <div className={styles.checkboxWrapper}>
                        <FormControl>
                          <Checkbox
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                        <div className={styles.checkboxLabel}>
                          <FormLabel>
                            {t("Profile", "Hidden") ?? "Hidden"}
                          </FormLabel>
                          <FormDescription>
                            {t("Profile", "Hide this link from your public profile.") ??
                              "Hide this link from your public profile."}
                          </FormDescription>
                        </div>
                      </div>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className={styles.actions}>
                  <Button type="button" variant="outline" onClick={handleCancel}>
                    {t("Profile", "Cancel") ?? "Cancel"}
                  </Button>
                  <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting
                      ? (t("Profile", "Saving...") ?? "Saving...")
                      : editingLink
                      ? (t("Profile", "Update Link") ?? "Update Link")
                      : (t("Profile", "Add Link") ?? "Add Link")}
                  </Button>
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
