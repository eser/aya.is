"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Input } from "@/shared/components/ui/input.tsx";
import { Label } from "@/shared/components/ui/label.tsx";
import { Textarea } from "@/shared/components/ui/textarea.tsx";
import { Badge } from "@/shared/components/ui/badge.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shared/components/ui/dialog.tsx";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/shared/components/ui/alert-dialog.tsx";
import styles from "./page.module.css";

type ProfilePageBrief = {
  id: string;
  slug: string;
  cover_picture_uri: string | null;
  title: string;
  summary: string;
};

type CreatePageForm = {
  slug: string;
  title: string;
  summary: string;
  content: string;
  cover_picture_uri: string | null;
};

export default function ProfilePagesPage() {
  const { t, locale } = useTranslations();
  const { isAuthenticated, isLoading } = useAuth();
  const params = useParams();
  const profileSlug = params.slug as string;

  const [pages, setPages] = useState<ProfilePageBrief[]>([]);
  const [isLoadingPages, setIsLoadingPages] = useState(true);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState<CreatePageForm>({
    slug: "",
    title: "",
    summary: "",
    content: "",
    cover_picture_uri: null,
  });

  const loadPages = async () => {
    if (!isAuthenticated || !profileSlug) return;

    try {
      setIsLoadingPages(true);
      const response = await backend.listProfilePages(locale.code, profileSlug);
      if (response !== null) {
        setPages(response);
      }
    } catch (err) {
      console.error("Failed to load pages:", err);
      setError(t("Profile", "Failed to load pages") ?? "Failed to load pages");
    } finally {
      setIsLoadingPages(false);
    }
  };

  useEffect(() => {
    loadPages();
  }, [isAuthenticated, profileSlug, locale]);

  const handleCreatePage = async () => {
    if (!createForm.slug || !createForm.title || !createForm.summary) {
      setError(t("Profile", "Please fill in all required fields") ?? "Please fill in all required fields");
      return;
    }

    try {
      setIsSubmitting(true);
      setError(null);

      await backend.createProfilePage(locale.code, profileSlug, {
        slug: createForm.slug,
        title: createForm.title,
        summary: createForm.summary,
        content: createForm.content,
        cover_picture_uri: createForm.cover_picture_uri,
        published_at: null,
      });

      // Reset form and close dialog
      setCreateForm({
        slug: "",
        title: "",
        summary: "",
        content: "",
        cover_picture_uri: null,
      });
      setIsCreateDialogOpen(false);

      // Reload pages
      await loadPages();
    } catch (err) {
      console.error("Failed to create page:", err);
      setError(t("Profile", "Failed to create page") ?? "Failed to create page");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeletePage = async (pageId: string) => {
    try {
      await backend.deleteProfilePage(locale.code, profileSlug, pageId);
      await loadPages();
    } catch (err) {
      console.error("Failed to delete page:", err);
      setError(t("Profile", "Failed to delete page") ?? "Failed to delete page");
    }
  };

  if (isLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>
          {t("Profile", "Loading...") ?? "Loading..."}
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className={styles.container}>
        <div className={styles.error}>
          {t("Profile", "Please log in to manage pages") ?? "Please log in to manage pages"}
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>
          {t("Profile", "Custom Pages") ?? "Custom Pages"}
        </h1>
        <p className={styles.description}>
          {t("Profile", "Create and manage custom content pages for your profile.") ??
            "Create and manage custom content pages for your profile."}
        </p>
      </div>

      {error && (
        <div className={styles.error}>
          {error}
        </div>
      )}

      <div className={styles.actions}>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>{t("Profile", "Create New Page") ?? "Create New Page"}</Button>
          </DialogTrigger>
          <DialogContent className={styles.dialog}>
            <DialogHeader>
              <DialogTitle>{t("Profile", "Create New Page") ?? "Create New Page"}</DialogTitle>
              <DialogDescription>
                {t("Profile", "Create a new custom page for your profile with rich content.") ??
                  "Create a new custom page for your profile with rich content."}
              </DialogDescription>
            </DialogHeader>
            <div className={styles.form}>
              <div className={styles.field}>
                <Label htmlFor="page-slug">{t("Profile", "Page Slug") ?? "Page Slug"} *</Label>
                <Input
                  id="page-slug"
                  value={createForm.slug}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, slug: e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "-") })}
                  placeholder={t("Profile", "my-custom-page") ?? "my-custom-page"}
                  className={styles.input}
                />
                <p className={styles.hint}>
                  {t("Profile", "URL-friendly identifier for your page") ?? "URL-friendly identifier for your page"}
                </p>
              </div>
              <div className={styles.field}>
                <Label htmlFor="page-title">{t("Profile", "Page Title") ?? "Page Title"} *</Label>
                <Input
                  id="page-title"
                  value={createForm.title}
                  onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
                  placeholder={t("Profile", "Enter page title") ?? "Enter page title"}
                  className={styles.input}
                />
              </div>
              <div className={styles.field}>
                <Label htmlFor="page-summary">{t("Profile", "Page Summary") ?? "Page Summary"} *</Label>
                <Textarea
                  id="page-summary"
                  value={createForm.summary}
                  onChange={(e) => setCreateForm({ ...createForm, summary: e.target.value })}
                  placeholder={t("Profile", "Brief description of your page") ?? "Brief description of your page"}
                  className={styles.textarea}
                  rows={3}
                />
              </div>
              <div className={styles.field}>
                <Label htmlFor="page-content">{t("Profile", "Page Content") ?? "Page Content"}</Label>
                <Textarea
                  id="page-content"
                  value={createForm.content}
                  onChange={(e) => setCreateForm({ ...createForm, content: e.target.value })}
                  placeholder={t("Profile", "Write your content here (supports Markdown)...") ??
                    "Write your content here (supports Markdown)..."}
                  className={styles.textarea}
                  rows={8}
                />
              </div>
              <div className={styles.field}>
                <Label htmlFor="page-cover">{t("Profile", "Cover Image URL") ?? "Cover Image URL"}</Label>
                <Input
                  id="page-cover"
                  value={createForm.cover_picture_uri ?? ""}
                  onChange={(e) => setCreateForm({ ...createForm, cover_picture_uri: e.target.value })}
                  placeholder={t("Profile", "https://example.com/image.jpg") ?? "https://example.com/image.jpg"}
                  className={styles.input}
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsCreateDialogOpen(false)}
                disabled={isSubmitting}
              >
                {t("Profile", "Cancel") ?? "Cancel"}
              </Button>
              <Button
                onClick={handleCreatePage}
                disabled={isSubmitting || !createForm.slug || !createForm.title || !createForm.summary}
              >
                {isSubmitting
                  ? (t("Profile", "Creating...") ?? "Creating...")
                  : (t("Profile", "Create Page") ?? "Create Page")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className={styles.pagesList}>
        {isLoadingPages
          ? (
            <div className={styles.loading}>
              {t("Profile", "Loading pages...") ?? "Loading pages..."}
            </div>
          )
          : pages.length === 0
          ? (
            <Card className={styles.emptyState}>
              <CardContent className={styles.emptyContent}>
                <span className={styles.emptyIcon}>📄</span>
                <h3 className={styles.emptyTitle}>
                  {t("Profile", "No custom pages yet") ?? "No custom pages yet"}
                </h3>
                <p className={styles.emptyDescription}>
                  {t("Profile", "Create your first custom page to showcase additional content.") ??
                    "Create your first custom page to showcase additional content."}
                </p>
              </CardContent>
            </Card>
          )
          : (
            pages.map((page) => (
              <Card key={page.id} className={styles.pageCard}>
                <CardHeader>
                  <div className={styles.pageHeader}>
                    <div className={styles.pageInfo}>
                      <CardTitle className={styles.pageTitle}>{page.title}</CardTitle>
                      <Badge variant="secondary" className={styles.pageSlug}>/{page.slug}</Badge>
                    </div>
                    <div className={styles.pageActions}>
                      <Button variant="outline" size="sm">
                        {t("Profile", "Edit") ?? "Edit"}
                      </Button>
                      <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button variant="outline" size="sm" className={styles.deleteButton}>
                            {t("Profile", "Delete") ?? "Delete"}
                          </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>
                              {t("Profile", "Delete Page") ?? "Delete Page"}
                            </AlertDialogTitle>
                            <AlertDialogDescription>
                              {t(
                                "Profile",
                                "Are you sure you want to delete this page? This action cannot be undone.",
                              ) ?? "Are you sure you want to delete this page? This action cannot be undone."}
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>
                              {t("Profile", "Cancel") ?? "Cancel"}
                            </AlertDialogCancel>
                            <AlertDialogAction
                              onClick={() => handleDeletePage(page.id)}
                              className={styles.deleteAction}
                            >
                              {t("Profile", "Delete") ?? "Delete"}
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <CardDescription className={styles.pageSummary}>
                    {page.summary}
                  </CardDescription>
                  {page.cover_picture_uri && (
                    <div className={styles.coverPreview}>
                      <img
                        src={page.cover_picture_uri}
                        alt={page.title}
                        className={styles.coverImage}
                      />
                    </div>
                  )}
                </CardContent>
              </Card>
            ))
          )}
      </div>
    </div>
  );
}
