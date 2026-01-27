import * as React from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import {
  AlertTriangle,
  ArrowLeft,
  Check,
  ImagePlus,
  Images,
  Info,
  Loader2,
  Megaphone,
  Newspaper,
  PanelLeftClose,
  PanelLeftOpen,
  PencilLine,
  Presentation,
  Upload,
} from "lucide-react";
import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger } from "@/components/ui/select";
import type { StoryKind } from "@/modules/backend/types";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { insertTextAtCursor, MarkdownEditor, wrapSelectedText } from "./markdown-editor";
import { PreviewPanel } from "./preview-panel";
import { EditorToolbar, type FormatAction, type ViewMode } from "./editor-toolbar";
import { type ContentStatus, EditorActions } from "./editor-actions";
import { ImageUploadModal } from "./image-upload-modal";
import styles from "./content-editor.module.css";
import { cn } from "@/lib/utils";
import { slugify } from "@/lib/slugify";
import { isAllowedURI } from "@/config";
import { backend } from "@/modules/backend/backend";
import { getEntityConfig } from "./entity-types";

type SlugAvailability = {
  isChecking: boolean;
  isAvailable: boolean | null;
  message: string | null;
  severity: "error" | "warning" | "" | null;
};

// Schema for optional URL validation (null, empty string, or valid http/https URL)
const optionalUrlSchema = z.union([
  z.literal(null),
  z.literal(""),
  z.string().url().refine(
    (url) => url.startsWith("http://") || url.startsWith("https://"),
    { message: "URL must use http or https protocol" },
  ),
]);

// Helper to format date as YYYYMMDD-
function formatDatePrefix(dateString: string | null): string | null {
  if (dateString === null || dateString === "") return null;
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return null;
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}${month}${day}-`;
}

// Validate slug starts with YYYYMMDD- of publish date
function validateSlugPrefix(
  slug: string,
  publishedAt: string | null,
): { valid: boolean; expectedPrefix: string | null } {
  const expectedPrefix = formatDatePrefix(publishedAt);
  if (expectedPrefix === null) {
    return { valid: true, expectedPrefix: null };
  }
  const valid = slug.startsWith(expectedPrefix);
  return { valid, expectedPrefix };
}

// Known API messages that need translation
const API_MESSAGE_KEYS: Record<string, string> = {
  "Slug must be at least 12 characters": "ContentEditor.Slug must be at least 12 characters",
  "This slug is already taken": "ContentEditor.This slug is already taken",
  "This slug was previously used": "ContentEditor.This slug was previously used",
  "This slug is reserved": "Profile.This slug is reserved",
};

export type ContentType = "story" | "page";

export type ContentEditorData = {
  title: string;
  slug: string;
  summary: string;
  content: string;
  storyPictureUri?: string | null;
  status: ContentStatus;
  publishedAt?: string | null;
  isFeatured?: boolean;
  kind?: StoryKind;
};

type ContentEditorProps = {
  locale: string;
  profileSlug: string;
  contentType: ContentType;
  initialData: ContentEditorData;
  backUrl: string;
  userKind?: string;
  validateSlugDatePrefix?: boolean;
  onSave: (data: ContentEditorData) => Promise<void>;
  onDelete?: () => Promise<void>;
  isNew?: boolean;
  excludeId?: string;
};

export function ContentEditor(props: ContentEditorProps) {
  const {
    locale,
    profileSlug,
    contentType,
    initialData,
    backUrl,
    userKind,
    validateSlugDatePrefix: shouldValidateSlugDatePrefix = false,
    onSave,
    onDelete,
    isNew = false,
    excludeId,
  } = props;

  const isAdmin = userKind === "admin";

  // Get entity-specific configuration
  const entityConfig = getEntityConfig(contentType);
  const imageFieldConfig = entityConfig.imageFields[0];

  const { t } = useTranslation();

  // Form state
  const [title, setTitle] = React.useState(initialData.title);
  const [slug, setSlug] = React.useState(initialData.slug);
  const [summary, setSummary] = React.useState(initialData.summary);
  const [content, setContent] = React.useState(initialData.content);
  const [storyPictureUri, setStoryPictureUri] = React.useState(
    initialData.storyPictureUri ?? null,
  );
  const [status, setStatus] = React.useState<ContentStatus>(initialData.status);
  const [publishedAt, setPublishedAt] = React.useState(
    initialData.publishedAt ?? null,
  );
  const [isFeatured, setIsFeatured] = React.useState(
    initialData.isFeatured ?? false,
  );
  const [kind, setKind] = React.useState<StoryKind>(
    initialData.kind ?? "article",
  );

  // UI state
  const [viewMode, setViewMode] = React.useState<ViewMode>("split");
  const [isSaving, setIsSaving] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);
  const [showImageModal, setShowImageModal] = React.useState(false);
  const [showStoryPictureModal, setShowStoryPictureModal] = React.useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = React.useState(false);
  const [slugError, setSlugError] = React.useState<string | null>(null);
  const [showSlugValidation, setShowSlugValidation] = React.useState(false);
  const [slugManuallyEdited, setSlugManuallyEdited] = React.useState(!isNew);
  const [titleError, setTitleError] = React.useState<string | null>(null);
  const [showTitleValidation, setTitleTouched] = React.useState(!isNew);
  const [storyPictureUriError, setStoryPictureUriError] = React.useState<string | null>(null);
  const [slugAvailability, setSlugAvailability] = React.useState<SlugAvailability>({
    isChecking: false,
    isAvailable: null,
    message: null,
    severity: null,
  });

  // Validate slug on change
  React.useEffect(() => {
    // Basic slug validation
    if (slug.length === 0) {
      setSlugError(showSlugValidation ? t("ContentEditor.Slug is required") : null);
      return;
    }
    if (slug.length < 2) {
      setSlugError(t("ContentEditor.Slug must be at least 2 characters"));
      return;
    }
    if (slug.length > 100) {
      setSlugError(t("ContentEditor.Slug must be at most 100 characters"));
      return;
    }
    if (!/^[a-z0-9-]+$/.test(slug)) {
      setSlugError(t("ContentEditor.Slug can only contain lowercase letters, numbers, and hyphens"));
      return;
    }

    // Date prefix validation for stories with global slugs
    if (shouldValidateSlugDatePrefix && status === "published" && publishedAt !== null) {
      const { valid, expectedPrefix } = validateSlugPrefix(slug, publishedAt);
      if (!valid && expectedPrefix !== null) {
        setSlugError(t("ContentEditor.Slug must start with") + ` ${expectedPrefix}`);
        return;
      }
    }

    setSlugError(null);
  }, [slug, showSlugValidation, publishedAt, status, shouldValidateSlugDatePrefix]);

  // Refs to capture current values without triggering effect re-runs
  const statusRef = React.useRef(status);
  const publishedAtRef = React.useRef(publishedAt);
  statusRef.current = status;
  publishedAtRef.current = publishedAt;

  // Debounced slug availability check
  React.useEffect(() => {
    // Only check if slug passes basic format validation (not date prefix)
    const hasBasicError = slug.length < 3 || slug.length > 100 || !/^[a-z0-9-]+$/.test(slug);
    if (hasBasicError) {
      setSlugAvailability({ isChecking: false, isAvailable: null, message: null, severity: null });
      return;
    }

    // Skip availability check if slug hasn't changed from initial (for editing)
    if (!isNew && slug === initialData.slug) {
      setSlugAvailability({ isChecking: false, isAvailable: true, message: null, severity: null });
      return;
    }

    // Skip availability check if slug is untouched for new content
    if (isNew && !showSlugValidation) {
      setSlugAvailability({ isChecking: false, isAvailable: null, message: null, severity: null });
      return;
    }

    setSlugAvailability((prev) => ({ ...prev, isChecking: true }));

    const timeoutId = setTimeout(async () => {
      try {
        let result: { available: boolean; message?: string; severity?: "error" | "warning" | "" } | null = null;

        if (contentType === "story") {
          result = await backend.checkStorySlug(locale, slug, {
            excludeId,
            status: statusRef.current,
            publishedAt: publishedAtRef.current,
          });
        } else {
          result = await backend.checkPageSlug(locale, profileSlug, slug, {
            excludeId,
          });
        }

        if (result !== null) {
          // Determine message based on severity and translate if needed
          let message: string | null = null;
          if (!result.available || result.severity === "warning") {
            const apiMessage = result.message ?? (result.available ? null : "This slug is already taken");
            // Translate known API messages
            if (apiMessage !== null && apiMessage in API_MESSAGE_KEYS) {
              message = t(API_MESSAGE_KEYS[apiMessage]);
            } else {
              message = apiMessage;
            }
          }

          setSlugAvailability({
            isChecking: false,
            isAvailable: result.available,
            message,
            severity: result.severity ?? null,
          });
        } else {
          setSlugAvailability({
            isChecking: false,
            isAvailable: null,
            message: null,
            severity: null,
          });
        }
      } catch {
        setSlugAvailability({
          isChecking: false,
          isAvailable: null,
          message: null,
          severity: null,
        });
      }
    }, 500);

    return () => {
      clearTimeout(timeoutId);
    };
  }, [slug, locale, profileSlug, contentType, excludeId, isNew, initialData.slug, showSlugValidation]);

  // Validate title on change
  React.useEffect(() => {
    if (title.length === 0) {
      setTitleError(showTitleValidation ? t("ContentEditor.Title is required") : null);
      return;
    }
    if (title.length > 200) {
      setTitleError(t("ContentEditor.Title must be at most 200 characters"));
      return;
    }
    setTitleError(null);
  }, [title, showTitleValidation]);

  // Validate story picture URI on change
  React.useEffect(() => {
    const result = optionalUrlSchema.safeParse(storyPictureUri);
    if (!result.success) {
      setStoryPictureUriError(t("ContentEditor.Invalid URI"));
      return;
    }

    // For non-admin users, validate URI prefix
    if (!isAdmin && storyPictureUri !== null && storyPictureUri !== "") {
      const prefixes = imageFieldConfig.allowedPrefixes;

      if (!isAllowedURI(storyPictureUri, prefixes)) {
        setStoryPictureUriError(
          t("ContentEditor.URI must start with allowed prefix") + `: ${prefixes.join(", ")}`,
        );
        return;
      }
    }

    setStoryPictureUriError(null);
  }, [storyPictureUri, t, isAdmin, imageFieldConfig.allowedPrefixes]);

  // Check if there are unsaved changes
  const hasChanges = React.useMemo(() => {
    return (
      title !== initialData.title ||
      slug !== initialData.slug ||
      summary !== initialData.summary ||
      content !== initialData.content ||
      storyPictureUri !== (initialData.storyPictureUri ?? null) ||
      status !== initialData.status ||
      publishedAt !== (initialData.publishedAt ?? null) ||
      isFeatured !== (initialData.isFeatured ?? false) ||
      kind !== (initialData.kind ?? "article")
    );
  }, [title, slug, summary, content, storyPictureUri, status, publishedAt, isFeatured, kind, initialData]);

  // Auto-generate slug from title for new content (called on title blur)
  const generateSlugFromTitle = React.useCallback(() => {
    // Only auto-generate when user hasn't manually edited the slug
    if (slugManuallyEdited || title.trim() === "") {
      return;
    }

    const generatedSlug = slugify(title);

    // Get the date prefix from current slug or generate new one
    const datePrefixPattern = /^(\d{8}-).*$/;
    const match = slug.match(datePrefixPattern);
    const datePrefix = match !== null ? match[1] : "";

    setSlug(datePrefix + generatedSlug);
    if (!showSlugValidation) setShowSlugValidation(true);
  }, [slugManuallyEdited, title, slug, showSlugValidation]);

  const getCurrentData = (): ContentEditorData => ({
    title,
    slug,
    summary,
    content,
    storyPictureUri,
    status,
    publishedAt,
    isFeatured,
    kind,
  });

  const handleSave = async () => {
    // Mark fields as touched to show any validation errors
    setShowSlugValidation(true);
    setTitleTouched(true);

    // Check for empty required fields
    if (slug.length === 0 || title.length === 0) {
      return;
    }

    // Check for validation errors - only block on errors, not warnings
    const hasSlugError = slugError !== null ||
      (slugAvailability.isAvailable === false && slugAvailability.severity === "error");
    if (hasSlugError || titleError !== null) {
      return;
    }

    // Validate slug prefix for published content (stories only)
    if (shouldValidateSlugDatePrefix && status === "published" && publishedAt !== null) {
      const { valid, expectedPrefix } = validateSlugPrefix(slug, publishedAt);
      if (!valid && expectedPrefix !== null) {
        setSlugError(t("ContentEditor.Slug must start with") + ` ${expectedPrefix}`);
        return;
      }
    }

    // Validate URI prefix for non-admin users
    if (!isAdmin && storyPictureUri !== null && storyPictureUri !== "") {
      const prefixes = imageFieldConfig.allowedPrefixes;

      if (!isAllowedURI(storyPictureUri, prefixes)) {
        setStoryPictureUriError(
          t("ContentEditor.URI must start with allowed prefix") + `: ${prefixes.join(", ")}`,
        );
        return;
      }
    }

    setIsSaving(true);
    try {
      await onSave(getCurrentData());
    } finally {
      setIsSaving(false);
    }
  };

  const handlePublish = async () => {
    // Mark fields as touched to show any validation errors
    setShowSlugValidation(true);
    setTitleTouched(true);

    // Check for empty required fields
    if (slug.length === 0 || title.length === 0) {
      return;
    }

    // Check for validation errors - only block on errors, not warnings
    const hasSlugError = slugError !== null ||
      (slugAvailability.isAvailable === false && slugAvailability.severity === "error");
    if (hasSlugError || titleError !== null) {
      return;
    }

    // Set current date/time as publish date if not set
    const effectivePublishedAt = publishedAt ?? new Date().toISOString().slice(0, 16);

    // Validate slug prefix (stories only)
    if (shouldValidateSlugDatePrefix) {
      const { valid, expectedPrefix } = validateSlugPrefix(slug, effectivePublishedAt);
      if (!valid && expectedPrefix !== null) {
        setSlugError(t("ContentEditor.Slug must start with") + ` ${expectedPrefix}`);
        return;
      }
    }

    // Validate URI prefix for non-admin users
    if (!isAdmin && storyPictureUri !== null && storyPictureUri !== "") {
      const prefixes = imageFieldConfig.allowedPrefixes;

      if (!isAllowedURI(storyPictureUri, prefixes)) {
        setStoryPictureUriError(
          t("ContentEditor.URI must start with allowed prefix") + `: ${prefixes.join(", ")}`,
        );
        return;
      }
    }

    setStatus("published");
    setPublishedAt(effectivePublishedAt);
    setIsSaving(true);
    try {
      await onSave({ ...getCurrentData(), status: "published", publishedAt: effectivePublishedAt });
    } finally {
      setIsSaving(false);
    }
  };

  const handleUnpublish = async () => {
    setStatus("draft");
    setIsSaving(true);
    try {
      await onSave({ ...getCurrentData(), status: "draft" });
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async () => {
    if (onDelete === undefined) return;
    setIsDeleting(true);
    try {
      await onDelete();
    } finally {
      setIsDeleting(false);
    }
  };

  const handleFormat = (action: FormatAction) => {
    const textarea = document.querySelector(
      `.${styles.markdownTextarea}`,
    ) as HTMLTextAreaElement | null;
    if (textarea === null) return;

    const formatMap: Record<
      FormatAction,
      { prefix: string; suffix: string } | { insert: string }
    > = {
      bold: { prefix: "**", suffix: "**" },
      italic: { prefix: "_", suffix: "_" },
      h2: { insert: "\n## " },
      h3: { insert: "\n### " },
      ul: { insert: "\n- " },
      ol: { insert: "\n1. " },
      link: { prefix: "[", suffix: "](url)" },
      code: { prefix: "`", suffix: "`" },
      quote: { insert: "\n> " },
    };

    const format = formatMap[action];
    if ("insert" in format) {
      insertTextAtCursor(textarea, format.insert, setContent);
    } else {
      wrapSelectedText(textarea, format.prefix, format.suffix, setContent);
    }
  };

  const handleImageInsert = (url: string, alt: string) => {
    const textarea = document.querySelector(
      `.${styles.markdownTextarea}`,
    ) as HTMLTextAreaElement | null;
    if (textarea === null) return;

    const markdown = `![${alt}](${url})`;
    insertTextAtCursor(textarea, markdown, setContent);
  };

  const handleStoryPictureInsert = (url: string, _alt: string) => {
    setStoryPictureUri(url);
  };

  return (
    <div className={styles.editorContainer}>
      {/* Header */}
      <div className={styles.editorHeader}>
        <div className="flex items-center gap-3">
          <Link
            to={backUrl}
            onClick={(e: React.MouseEvent) => {
              e.preventDefault();
              globalThis.history.back();
            }}
          >
            <Button variant="outline" size="icon" className="rounded-full">
              <ArrowLeft className="size-4" />
            </Button>
          </Link>
          <h1 className="text-lg font-semibold">
            {isNew
              ? t(contentType === "story" ? "ContentEditor.New Story" : "ContentEditor.New Page")
              : t(contentType === "story" ? "ContentEditor.Edit Story" : "ContentEditor.Edit Page")}
          </h1>
        </div>

        <EditorActions
          status={status}
          isSaving={isSaving}
          isDeleting={isDeleting}
          hasChanges={hasChanges}
          onSave={handleSave}
          onPublish={handlePublish}
          onUnpublish={handleUnpublish}
          onDelete={handleDelete}
          canDelete={!isNew && onDelete !== undefined && isAdmin}
        />
      </div>

      {/* Main Content */}
      <div className={styles.editorMain}>
        {/* Sidebar - Metadata */}
        <div
          className={cn(
            styles.editorSidebar,
            sidebarCollapsed && styles.editorSidebarCollapsed,
          )}
        >
          <div className={styles.sidebarHeader}>
            <span className={cn(styles.sidebarTitle, sidebarCollapsed && "hidden")}>
              {t("ContentEditor.Metadata")}
            </span>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              title={sidebarCollapsed ? t("ContentEditor.Expand sidebar") : t("ContentEditor.Collapse sidebar")}
            >
              {sidebarCollapsed ? <PanelLeftOpen className="size-4" /> : <PanelLeftClose className="size-4" />}
            </Button>
          </div>

          {!sidebarCollapsed && (
            <div className={styles.metadataForm}>
              {/* Kind */}
              {contentType === "story" && (
                <Field className={styles.metadataField}>
                  <FieldLabel htmlFor="kind" className={styles.metadataLabel}>
                    {t("ContentEditor.Published At")}
                  </FieldLabel>
                  <Select value={kind} onValueChange={(value) => setKind(value as StoryKind)}>
                    <SelectTrigger id="kind">
                      <span className="flex items-center gap-2">
                        {kind === "article" && <PencilLine className="size-4" />}
                        {kind === "announcement" && <Megaphone className="size-4" />}
                        {kind === "news" && <Newspaper className="size-4" />}
                        {kind === "status" && <Info className="size-4" />}
                        {kind === "content" && <Images className="size-4" />}
                        {kind === "presentation" && <Presentation className="size-4" />}
                        {kind === "article" && t("Stories.Article")}
                        {kind === "announcement" && t("Stories.Announcement")}
                        {kind === "news" && t("ContentEditor.News")}
                        {kind === "status" && t("Stories.Status")}
                        {kind === "content" && t("Stories.Content")}
                        {kind === "presentation" && t("Stories.Presentation")}
                      </span>
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="article">
                        <span className="flex items-center gap-2">
                          <PencilLine className="size-4" />
                          {t("Stories.Article")}
                        </span>
                      </SelectItem>
                      <SelectItem value="announcement">
                        <span className="flex items-center gap-2">
                          <Megaphone className="size-4" />
                          {t("Stories.Announcement")}
                        </span>
                      </SelectItem>
                      <SelectItem value="news">
                        <span className="flex items-center gap-2">
                          <Newspaper className="size-4" />
                          {t("ContentEditor.News")}
                        </span>
                      </SelectItem>
                      <SelectItem value="status">
                        <span className="flex items-center gap-2">
                          <Info className="size-4" />
                          {t("Stories.Status")}
                        </span>
                      </SelectItem>
                      <SelectItem value="content">
                        <span className="flex items-center gap-2">
                          <Images className="size-4" />
                          {t("Stories.Content")}
                        </span>
                      </SelectItem>
                      <SelectItem value="presentation">
                        <span className="flex items-center gap-2">
                          <Presentation className="size-4" />
                          {t("Stories.Presentation")}
                        </span>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </Field>
              )}

              {/* Published At - visible for published content, editable only by admin */}
              {(status === "published" || isAdmin) && (
                <Field className={styles.metadataField}>
                  <FieldLabel htmlFor="published-at" className={styles.metadataLabel}>
                    {t("ContentEditor.Published At")}
                  </FieldLabel>
                  <Input
                    id="published-at"
                    type="text"
                    value={publishedAt ?? ""}
                    onChange={(e) => setPublishedAt(e.target.value || null)}
                    disabled={!isAdmin}
                  />
                </Field>
              )}

              {/* Is Featured - Admin only */}
              {isAdmin && contentType === "story" && (
                <Field className={styles.metadataField} orientation="horizontal">
                  <FieldLabel htmlFor="is-featured" className={styles.metadataLabel}>
                    {t("ContentEditor.Featured")}
                  </FieldLabel>
                  <Switch
                    id="is-featured"
                    checked={isFeatured}
                    onCheckedChange={setIsFeatured}
                  />
                </Field>
              )}

              {/* Title */}
              <Field className={styles.metadataField} data-invalid={titleError !== null || undefined}>
                <FieldLabel htmlFor="title" className={styles.metadataLabel}>
                  {t("ContentEditor.Title")}
                </FieldLabel>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => {
                    setTitle(e.target.value);
                    if (!showTitleValidation) setTitleTouched(true);
                  }}
                  onBlur={() => {
                    setTitleTouched(true);
                    generateSlugFromTitle();
                  }}
                  placeholder={t("ContentEditor.Enter title...")}
                  aria-invalid={titleError !== null || undefined}
                />
                {titleError !== null && <FieldError>{titleError}</FieldError>}
              </Field>

              {/* Summary */}
              <Field className={styles.metadataField}>
                <FieldLabel htmlFor="summary" className={styles.metadataLabel}>
                  {t("ContentEditor.Summary")}
                </FieldLabel>
                <Textarea
                  id="summary"
                  value={summary}
                  onChange={(e) => setSummary(e.target.value)}
                  placeholder={t("ContentEditor.Brief summary...")}
                  className="min-h-[80px]"
                />
              </Field>

              {/* Slug */}
              <Field
                className={styles.metadataField}
                data-invalid={(showSlugValidation && slugError !== null) ||
                  (!slugAvailability.isChecking && slugAvailability.severity === "error") || undefined}
              >
                <FieldLabel htmlFor="slug" className={styles.metadataLabel}>
                  {t("ContentEditor.Slug")}
                </FieldLabel>
                <div className="relative">
                  <Input
                    id="slug"
                    value={slug}
                    onChange={(e) => {
                      setSlug(e.target.value);
                      if (!showSlugValidation) setShowSlugValidation(true);
                      if (!slugManuallyEdited) setSlugManuallyEdited(true);
                    }}
                    onBlur={() => {
                      if (!showSlugValidation) setShowSlugValidation(true);
                    }}
                    placeholder={t("ContentEditor.url-friendly-slug")}
                    aria-invalid={(showSlugValidation && slugError !== null) ||
                      (!slugAvailability.isChecking && slugAvailability.severity === "error") || undefined}
                    className="pr-8"
                  />
                  {showSlugValidation && slug.length >= 3 && (
                    <div className="absolute right-2 top-1/2 -translate-y-1/2">
                      {slugAvailability.isChecking && <Loader2 className="size-4 animate-spin text-muted-foreground" />}
                      {!slugAvailability.isChecking && slugAvailability.isAvailable === true &&
                        slugAvailability.severity !== "warning" && <Check className="size-4 text-green-600" />}
                      {!slugAvailability.isChecking && slugAvailability.severity === "warning" && (
                        <AlertTriangle className="size-4 text-amber-500" />
                      )}
                    </div>
                  )}
                </div>
                {showSlugValidation && slugError !== null && <FieldError>{slugError}</FieldError>}
                {showSlugValidation && slugError === null && !slugAvailability.isChecking &&
                  slugAvailability.severity === "error" && slugAvailability.message !== null && (
                  <FieldError>{slugAvailability.message}</FieldError>
                )}
                {showSlugValidation && slugError === null && !slugAvailability.isChecking &&
                  slugAvailability.severity === "warning" && slugAvailability.message !== null && (
                  <p className="text-sm text-amber-600 mt-1">{slugAvailability.message}</p>
                )}
              </Field>

              {/* Picture URI (Story Picture or Cover Picture depending on entity type) */}
              <Field className={styles.metadataField} data-invalid={storyPictureUriError !== null || undefined}>
                <FieldLabel htmlFor="story-picture-uri" className={styles.metadataLabel}>
                  {t(imageFieldConfig.labelKey)}
                </FieldLabel>
                <div className="flex gap-2">
                  <Input
                    id="story-picture-uri"
                    value={storyPictureUri ?? ""}
                    onChange={(e) => setStoryPictureUri(e.target.value || null)}
                    placeholder="https://..."
                    className="flex-1"
                    aria-invalid={storyPictureUriError !== null || undefined}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="icon"
                    onClick={() => setShowStoryPictureModal(true)}
                    title={t("ContentEditor.Upload")}
                  >
                    <Upload className="size-4" />
                  </Button>
                  {contentType === "story" && !isNew && initialData.slug !== "" && (
                    <Link
                      to="/$locale/stories/$storyslug/cover"
                      params={{ locale, storyslug: initialData.slug }}
                    >
                      <Button
                        type="button"
                        variant="outline"
                        size="icon"
                        title={t("CoverDesigner.Design Cover")}
                      >
                        <ImagePlus className="size-4" />
                      </Button>
                    </Link>
                  )}
                </div>
                {storyPictureUriError !== null && <FieldError>{storyPictureUriError}</FieldError>}
                {storyPictureUri !== null && storyPictureUri !== "" && storyPictureUriError === null && (
                  <img
                    src={storyPictureUri}
                    alt="Picture preview"
                    className="mt-2 rounded-md max-h-32 w-full object-cover"
                  />
                )}
              </Field>
            </div>
          )}
        </div>

        {/* Editor Content */}
        <div className={styles.editorContent}>
          <EditorToolbar
            viewMode={viewMode}
            onViewModeChange={setViewMode}
            onFormat={handleFormat}
            onImageUpload={() => setShowImageModal(true)}
          />

          <div className={styles.editorPanels}>
            {/* Split View with Resizable Panels */}
            {viewMode === "split" && (
              <ResizablePanelGroup direction="horizontal" className="h-full">
                <ResizablePanel defaultSize={50} minSize={25}>
                  <div className={styles.editorPanel}>
                    <MarkdownEditor
                      value={content}
                      onChange={setContent}
                      placeholder={t("ContentEditor.Write your content in markdown...")}
                    />
                  </div>
                </ResizablePanel>
                <ResizableHandle withHandle />
                <ResizablePanel defaultSize={50} minSize={25}>
                  <div className={styles.editorPanel}>
                    <PreviewPanel content={content} />
                  </div>
                </ResizablePanel>
              </ResizablePanelGroup>
            )}

            {/* Editor Only */}
            {viewMode === "editor" && (
              <div className={styles.editorPanel}>
                <MarkdownEditor
                  value={content}
                  onChange={setContent}
                  placeholder={t("ContentEditor.Write your content in markdown...")}
                />
              </div>
            )}

            {/* Preview Only */}
            {viewMode === "preview" && (
              <div className={styles.editorPanel}>
                <PreviewPanel content={content} />
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Image Upload Modal for content */}
      <ImageUploadModal
        open={showImageModal}
        onOpenChange={setShowImageModal}
        onImageInsert={handleImageInsert}
        locale={locale}
        purpose="content-image"
      />

      {/* Image Upload Modal for story picture */}
      <ImageUploadModal
        open={showStoryPictureModal}
        onOpenChange={setShowStoryPictureModal}
        onImageInsert={handleStoryPictureInsert}
        locale={locale}
        purpose="cover-image"
      />
    </div>
  );
}

export type { ContentStatus };
