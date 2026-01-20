import * as React from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import {
  ArrowLeft,
  PanelLeftClose,
  PanelLeftOpen,
  Newspaper,
  PencilLine,
  Megaphone,
  Info,
  Images,
  Presentation,
  Upload,
} from "lucide-react";
import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { StoryKind } from "@/modules/backend/types";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";
import { MarkdownEditor, wrapSelectedText, insertTextAtCursor } from "./markdown-editor";
import { PreviewPanel } from "./preview-panel";
import {
  EditorToolbar,
  type ViewMode,
  type FormatAction,
} from "./editor-toolbar";
import { EditorActions, type ContentStatus } from "./editor-actions";
import { ImageUploadModal } from "./image-upload-modal";
import styles from "./content-editor.module.css";
import { cn } from "@/lib/utils";
import { isAllowedURI } from "@/config";
import { getEntityConfig } from "./entity-types";

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
function validateSlugPrefix(slug: string, publishedAt: string | null): { valid: boolean; expectedPrefix: string | null } {
  const expectedPrefix = formatDatePrefix(publishedAt);
  if (expectedPrefix === null) {
    return { valid: true, expectedPrefix: null };
  }
  const valid = slug.startsWith(expectedPrefix);
  return { valid, expectedPrefix };
}

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
};

export function ContentEditor(props: ContentEditorProps) {
  const {
    locale,
    contentType,
    initialData,
    backUrl,
    userKind,
    validateSlugDatePrefix: shouldValidateSlugDatePrefix = false,
    onSave,
    onDelete,
    isNew = false,
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
  const [storyPictureUriError, setStoryPictureUriError] = React.useState<string | null>(null);

  // Validate slug on change (only for stories with global slugs)
  React.useEffect(() => {
    if (shouldValidateSlugDatePrefix && status === "published" && publishedAt !== null) {
      const { valid, expectedPrefix } = validateSlugPrefix(slug, publishedAt);
      if (!valid && expectedPrefix !== null) {
        setSlugError(t("Editor.Slug must start with") + ` ${expectedPrefix}`);
      } else {
        setSlugError(null);
      }
    } else {
      setSlugError(null);
    }
  }, [slug, publishedAt, status, shouldValidateSlugDatePrefix, t]);

  // Validate story picture URI on change
  React.useEffect(() => {
    const result = optionalUrlSchema.safeParse(storyPictureUri);
    if (!result.success) {
      setStoryPictureUriError(t("Editor.Invalid URI"));
      return;
    }

    // For non-admin users, validate URI prefix
    if (!isAdmin && storyPictureUri !== null && storyPictureUri !== "") {
      const prefixes = imageFieldConfig.allowedPrefixes;

      if (!isAllowedURI(storyPictureUri, prefixes)) {
        setStoryPictureUriError(
          t("Editor.URI must start with allowed prefix") + `: ${prefixes.join(", ")}`,
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

  // Auto-generate slug from title for new content
  React.useEffect(() => {
    if (isNew && title !== "" && slug === "") {
      const generatedSlug = title
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, "")
        .replace(/\s+/g, "-")
        .replace(/-+/g, "-")
        .trim();
      setSlug(generatedSlug);
    }
  }, [title, slug, isNew]);

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
    // Validate slug prefix for published content (stories only)
    if (shouldValidateSlugDatePrefix && status === "published" && publishedAt !== null) {
      const { valid, expectedPrefix } = validateSlugPrefix(slug, publishedAt);
      if (!valid && expectedPrefix !== null) {
        setSlugError(t("Editor.Slug must start with") + ` ${expectedPrefix}`);
        return;
      }
    }

    // Validate URI prefix for non-admin users
    if (!isAdmin && storyPictureUri !== null && storyPictureUri !== "") {
      const prefixes = imageFieldConfig.allowedPrefixes;

      if (!isAllowedURI(storyPictureUri, prefixes)) {
        setStoryPictureUriError(
          t("Editor.URI must start with allowed prefix") + `: ${prefixes.join(", ")}`,
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
    // Set current date/time as publish date if not set
    const effectivePublishedAt = publishedAt ?? new Date().toISOString().slice(0, 16);

    // Validate slug prefix (stories only)
    if (shouldValidateSlugDatePrefix) {
      const { valid, expectedPrefix } = validateSlugPrefix(slug, effectivePublishedAt);
      if (!valid && expectedPrefix !== null) {
        setSlugError(t("Editor.Slug must start with") + ` ${expectedPrefix}`);
        return;
      }
    }

    // Validate URI prefix for non-admin users
    if (!isAdmin && storyPictureUri !== null && storyPictureUri !== "") {
      const prefixes = imageFieldConfig.allowedPrefixes;

      if (!isAllowedURI(storyPictureUri, prefixes)) {
        setStoryPictureUriError(
          t("Editor.URI must start with allowed prefix") + `: ${prefixes.join(", ")}`,
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
              ? t(contentType === "story" ? "Editor.New Story" : "Editor.New Page")
              : t(contentType === "story" ? "Editor.Edit Story" : "Editor.Edit Page")}
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
              {t("Editor.Metadata")}
            </span>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              title={sidebarCollapsed ? t("Editor.Expand sidebar") : t("Editor.Collapse sidebar")}
            >
              {sidebarCollapsed ? (
                <PanelLeftOpen className="size-4" />
              ) : (
                <PanelLeftClose className="size-4" />
              )}
            </Button>
          </div>

          {!sidebarCollapsed && (
            <div className={styles.metadataForm}>
              {/* Kind */}
              {contentType === "story" && (
                <div className={styles.metadataField}>
                  <Label htmlFor="kind" className={styles.metadataLabel}>
                    {t("Editor.Kind")}
                  </Label>
                  <Select value={kind} onValueChange={(value) => setKind(value as StoryKind)}>
                    <SelectTrigger id="kind">
                      <SelectValue />
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
                          {t("Editor.News")}
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
                </div>
              )}

              {/* Slug */}
              <div className={styles.metadataField}>
                <Label htmlFor="slug" className={styles.metadataLabel}>
                  {t("Editor.Slug")}
                </Label>
                <Input
                  id="slug"
                  value={slug}
                  onChange={(e) => setSlug(e.target.value)}
                  placeholder={t("Editor.url-friendly-slug")}
                  className={slugError !== null ? "border-destructive" : ""}
                />
                {slugError !== null && (
                  <p className="text-sm text-destructive mt-1">{slugError}</p>
                )}
              </div>

              {/* Published At - visible for published content, editable only by admin */}
              {(status === "published" || isAdmin) && (
                <div className={styles.metadataField}>
                  <Label htmlFor="published-at" className={styles.metadataLabel}>
                    {t("Editor.Published At")}
                  </Label>
                  <Input
                    id="published-at"
                    type="text"
                    value={publishedAt ?? ""}
                    onChange={(e) =>
                      setPublishedAt(e.target.value || null)
                    }
                    disabled={!isAdmin}
                  />
                </div>
              )}

              {/* Is Featured - Admin only */}
              {isAdmin && contentType === "story" && (
                <div className={styles.metadataField}>
                  <div className="flex items-center justify-between">
                    <Label htmlFor="is-featured" className={styles.metadataLabel}>
                      {t("Editor.Featured")}
                    </Label>
                    <Switch
                      id="is-featured"
                      checked={isFeatured}
                      onCheckedChange={setIsFeatured}
                    />
                  </div>
                </div>
              )}

              {/* Title */}
              <div className={styles.metadataField}>
                <Label htmlFor="title" className={styles.metadataLabel}>
                  {t("Editor.Title")}
                </Label>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder={t("Editor.Enter title...")}
                />
              </div>

              {/* Summary */}
              <div className={styles.metadataField}>
                <Label htmlFor="summary" className={styles.metadataLabel}>
                  {t("Editor.Summary")}
                </Label>
                <Textarea
                  id="summary"
                  value={summary}
                  onChange={(e) => setSummary(e.target.value)}
                  placeholder={t("Editor.Brief summary...")}
                  className="min-h-[80px]"
                />
              </div>

              {/* Picture URI (Story Picture or Cover Picture depending on entity type) */}
              <div className={styles.metadataField}>
                <Label htmlFor="story-picture-uri" className={styles.metadataLabel}>
                  {t(imageFieldConfig.labelKey)}
                </Label>
                <div className="flex gap-2">
                  <Input
                    id="story-picture-uri"
                    value={storyPictureUri ?? ""}
                    onChange={(e) =>
                      setStoryPictureUri(e.target.value || null)
                    }
                    placeholder="https://..."
                    className={cn("flex-1", storyPictureUriError !== null && "border-destructive")}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="icon"
                    onClick={() => setShowStoryPictureModal(true)}
                    title={t("Editor.Upload")}
                  >
                    <Upload className="size-4" />
                  </Button>
                </div>
                {storyPictureUriError !== null && (
                  <p className="text-sm text-destructive mt-1">{storyPictureUriError}</p>
                )}
                {storyPictureUri !== null && storyPictureUri !== "" && storyPictureUriError === null && (
                  <img
                    src={storyPictureUri}
                    alt="Picture preview"
                    className="mt-2 rounded-md max-h-32 w-full object-cover"
                  />
                )}
              </div>
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
                      placeholder={t("Editor.Write your content in markdown...")}
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
                  placeholder={t("Editor.Write your content in markdown...")}
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
