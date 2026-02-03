// Cover Generator Component
// Main reusable component for generating story covers

import * as React from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import type { Story } from "@/modules/backend/types.ts";
import type { StoryData, CoverOptions, TemplateId } from "@/lib/cover-generator/types.ts";
import { defaultCoverOptions } from "@/lib/cover-generator/types.ts";
import { getTemplate } from "@/lib/cover-generator/templates.ts";
import { downloadCanvas } from "@/lib/cover-generator/canvas-renderer.ts";
import { renderCoverSvg, downloadSvg } from "@/lib/cover-generator/svg-renderer.ts";
import { uploadCoverImage, loadImage } from "@/lib/cover-generator/upload.ts";
import { backend } from "@/modules/backend/backend.ts";
import { Button } from "@/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Download, FileCode, FileImage, ImagePlus, Loader2, ArrowLeft, ChevronDown } from "lucide-react";
import { CanvasPreview } from "./canvas-preview.tsx";
import { TemplateSelector } from "./template-selector.tsx";
import { CustomizationPanel } from "./customization-panel.tsx";
import styles from "./cover-generator.module.css";

interface CoverGeneratorProps {
  story: Story;
  locale: string;
  onBack?: () => void;
  onCoverSet?: (coverUrl: string) => void;
}

export function CoverGenerator(props: CoverGeneratorProps) {
  const { t } = useTranslation();
  const { story, locale } = props;
  const canvasRef = React.useRef<HTMLCanvasElement | null>(null);
  const [isUploading, setIsUploading] = React.useState(false);
  const [authorImageDataUrl, setAuthorImageDataUrl] = React.useState<string | null>(null);

  // Load author image as data URL for SVG export
  React.useEffect(() => {
    const loadAuthorImage = async () => {
      if (story.author_profile?.profile_picture_uri === null || story.author_profile?.profile_picture_uri === undefined) {
        setAuthorImageDataUrl(null);
        return;
      }

      try {
        const img = await loadImage(story.author_profile.profile_picture_uri);
        const canvas = document.createElement("canvas");
        canvas.width = img.width;
        canvas.height = img.height;
        const ctx = canvas.getContext("2d");
        if (ctx !== null) {
          ctx.drawImage(img, 0, 0);
          setAuthorImageDataUrl(canvas.toDataURL("image/png"));
        }
      } catch {
        setAuthorImageDataUrl(null);
      }
    };

    loadAuthorImage();
  }, [story.author_profile?.profile_picture_uri]);

  // Get localized kind label
  const getKindLabel = (kind: string): string => {
    const kindKey = kind.charAt(0).toUpperCase() + kind.slice(1);
    return t(`Stories.${kindKey}`, kindKey);
  };

  // Convert Story to StoryData
  const storyData: StoryData = React.useMemo(() => ({
    title: story.title ?? "Untitled",
    summary: story.summary,
    kind: story.kind,
    kindLabel: getKindLabel(story.kind),
    authorName: story.author_profile?.title ?? null,
    authorAvatarUrl: story.author_profile?.profile_picture_uri ?? null,
    publishedAt: story.created_at,
  }), [story, t]);

  // Initialize options with template defaults
  const [options, setOptions] = React.useState<CoverOptions>(() => ({
    ...defaultCoverOptions,
    ...getTemplate("classic").defaults,
    locale,
  }));

  // Handle template change
  const handleTemplateChange = (templateId: TemplateId) => {
    const template = getTemplate(templateId);
    setOptions((prev) => ({
      ...prev,
      ...template.defaults,
      templateId,
      // Keep user's custom text overrides and locale
      locale: prev.locale,
      titleOverride: prev.titleOverride,
      subtitleOverride: prev.subtitleOverride,
    }));
  };

  // Handle options change
  const handleOptionsChange = (partialOptions: Partial<CoverOptions>) => {
    setOptions((prev) => ({ ...prev, ...partialOptions }));
  };

  // Handle download PNG
  const handleDownloadPng = () => {
    if (canvasRef.current === null) {
      toast.error(t("CoverDesigner.Failed to generate image"));
      return;
    }

    const filename = `${story.slug ?? "cover"}-cover.png`;
    downloadCanvas(canvasRef.current, filename);
    toast.success(t("CoverDesigner.Cover downloaded"));
  };

  // Handle download SVG
  const handleDownloadSvg = () => {
    const svgContent = renderCoverSvg(storyData, options, authorImageDataUrl);
    const filename = `${story.slug ?? "cover"}-cover.svg`;
    downloadSvg(svgContent, filename);
    toast.success(t("CoverDesigner.Cover downloaded"));
  };

  // Handle set as cover
  const handleSetAsCover = async () => {
    if (canvasRef.current === null) {
      toast.error(t("CoverDesigner.Failed to generate image"));
      return;
    }

    if (story.author_profile === null) {
      toast.error(t("CoverDesigner.Story has no author"));
      return;
    }

    setIsUploading(true);

    try {
      // Upload the cover image
      const result = await uploadCoverImage(
        canvasRef.current,
        locale,
        story.slug ?? story.id,
      );

      if (!result.success || result.publicUrl === undefined) {
        toast.error(result.error ?? t("CoverDesigner.Failed to upload image"));
        return;
      }

      // Update story with new cover
      // Note: backend.updateStory returns null on success when no data is returned
      // The fetcher throws on error, so reaching here means success
      await backend.updateStory(
        locale,
        story.author_profile.slug,
        story.id,
        {
          slug: story.slug ?? "",
          story_picture_uri: result.publicUrl,
        },
      );

      toast.success(t("CoverDesigner.Cover set successfully"));
      props.onCoverSet?.(result.publicUrl);
    } catch (error) {
      const message = error instanceof Error ? error.message : t("CoverDesigner.An error occurred");
      toast.error(message);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className="flex items-center gap-3">
          {props.onBack !== undefined && (
            <Button variant="outline" size="icon" className="rounded-full" onClick={props.onBack}>
              <ArrowLeft className="size-4" />
            </Button>
          )}
          <h1 className={styles.title}>{t("CoverDesigner.Design Cover")}</h1>
        </div>
        <div className={styles.headerActions}>
          <DropdownMenu>
            <DropdownMenuTrigger className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2">
              <Download className="size-4" />
              {t("CoverDesigner.Download")}
              <ChevronDown className="size-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={handleDownloadPng}>
                <FileImage className="size-4" />
                {t("CoverDesigner.Download PNG")}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleDownloadSvg}>
                <FileCode className="size-4" />
                {t("CoverDesigner.Download SVG")}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Button onClick={handleSetAsCover} disabled={isUploading}>
            {isUploading ? (
              <Loader2 className="mr-2 size-4 animate-spin" />
            ) : (
              <ImagePlus className="mr-2 size-4" />
            )}
            {t("CoverDesigner.Set as Cover")}
          </Button>
        </div>
      </div>

      {/* Main content */}
      <div className={styles.content}>
        {/* Preview */}
        <div className={styles.previewSection}>
          <CanvasPreview
            story={storyData}
            options={options}
            canvasRef={canvasRef}
          />
        </div>

        {/* Controls */}
        <div className={styles.controlsSection}>
          <TemplateSelector
            selectedTemplate={options.templateId}
            onSelect={handleTemplateChange}
          />
          <CustomizationPanel
            options={options}
            onChange={handleOptionsChange}
          />
        </div>
      </div>
    </div>
  );
}
