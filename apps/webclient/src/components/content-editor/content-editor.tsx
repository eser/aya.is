import * as React from "react";
import { ArrowLeft } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
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

export type ContentType = "story" | "page";

export type ContentEditorData = {
  title: string;
  slug: string;
  summary: string;
  content: string;
  coverImageUrl?: string | null;
  status: ContentStatus;
};

type ContentEditorProps = {
  locale: string;
  profileSlug: string;
  contentType: ContentType;
  initialData: ContentEditorData;
  backUrl: string;
  backLabel?: string;
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
    backLabel = "Back",
    onSave,
    onDelete,
    isNew = false,
  } = props;

  // Form state
  const [title, setTitle] = React.useState(initialData.title);
  const [slug, setSlug] = React.useState(initialData.slug);
  const [summary, setSummary] = React.useState(initialData.summary);
  const [content, setContent] = React.useState(initialData.content);
  const [coverImageUrl, setCoverImageUrl] = React.useState(
    initialData.coverImageUrl ?? null,
  );
  const [status, setStatus] = React.useState<ContentStatus>(initialData.status);

  // UI state
  const [viewMode, setViewMode] = React.useState<ViewMode>("split");
  const [isSaving, setIsSaving] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);
  const [showImageModal, setShowImageModal] = React.useState(false);

  // Check if there are unsaved changes
  const hasChanges = React.useMemo(() => {
    return (
      title !== initialData.title ||
      slug !== initialData.slug ||
      summary !== initialData.summary ||
      content !== initialData.content ||
      coverImageUrl !== (initialData.coverImageUrl ?? null) ||
      status !== initialData.status
    );
  }, [title, slug, summary, content, coverImageUrl, status, initialData]);

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
    coverImageUrl,
    status,
  });

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await onSave(getCurrentData());
    } finally {
      setIsSaving(false);
    }
  };

  const handlePublish = async () => {
    setStatus("published");
    setIsSaving(true);
    try {
      await onSave({ ...getCurrentData(), status: "published" });
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

  return (
    <div className={styles.editorContainer}>
      {/* Header */}
      <div className={styles.editorHeader}>
        <div className="flex items-center gap-3">
          <Link to={backUrl}>
            <Button variant="ghost" size="sm">
              <ArrowLeft className="mr-1.5 size-4" />
              {backLabel}
            </Button>
          </Link>
          <h1 className="text-lg font-semibold">
            {isNew
              ? `New ${contentType === "story" ? "Story" : "Page"}`
              : `Edit ${contentType === "story" ? "Story" : "Page"}`}
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
          showDelete={!isNew && onDelete !== undefined}
        />
      </div>

      {/* Main Content */}
      <div className={styles.editorMain}>
        {/* Sidebar - Metadata */}
        <div className={styles.editorSidebar}>
          <div className={styles.metadataForm}>
            <div className={styles.metadataField}>
              <Label htmlFor="title" className={styles.metadataLabel}>
                Title
              </Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter title..."
              />
            </div>

            <div className={styles.metadataField}>
              <Label htmlFor="slug" className={styles.metadataLabel}>
                Slug
              </Label>
              <Input
                id="slug"
                value={slug}
                onChange={(e) => setSlug(e.target.value)}
                placeholder="url-friendly-slug"
              />
            </div>

            <div className={styles.metadataField}>
              <Label htmlFor="summary" className={styles.metadataLabel}>
                Summary
              </Label>
              <Textarea
                id="summary"
                value={summary}
                onChange={(e) => setSummary(e.target.value)}
                placeholder="Brief summary..."
                className="min-h-[80px]"
              />
            </div>

            <div className={styles.metadataField}>
              <Label htmlFor="cover-image" className={styles.metadataLabel}>
                Cover Image URL
              </Label>
              <Input
                id="cover-image"
                value={coverImageUrl ?? ""}
                onChange={(e) =>
                  setCoverImageUrl(e.target.value || null)
                }
                placeholder="https://..."
              />
              {coverImageUrl !== null && coverImageUrl !== "" && (
                <img
                  src={coverImageUrl}
                  alt="Cover preview"
                  className="mt-2 rounded-md max-h-32 w-full object-cover"
                />
              )}
            </div>
          </div>
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
            {/* Editor Panel */}
            {(viewMode === "editor" || viewMode === "split") && (
              <div
                className={cn(
                  styles.editorPanel,
                  viewMode === "split" && styles.editorPanelSplit,
                )}
              >
                <MarkdownEditor
                  value={content}
                  onChange={setContent}
                  placeholder="Write your content in markdown..."
                />
              </div>
            )}

            {/* Divider */}
            {viewMode === "split" && (
              <div className={styles.editorPanelDivider} />
            )}

            {/* Preview Panel */}
            {(viewMode === "preview" || viewMode === "split") && (
              <div
                className={cn(
                  styles.editorPanel,
                  viewMode === "split" && styles.editorPanelSplit,
                )}
              >
                <PreviewPanel content={content} />
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Image Upload Modal */}
      <ImageUploadModal
        open={showImageModal}
        onOpenChange={setShowImageModal}
        onImageInsert={handleImageInsert}
        locale={locale}
        purpose="content-image"
      />
    </div>
  );
}

export type { ContentStatus };
