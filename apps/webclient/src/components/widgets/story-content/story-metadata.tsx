import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import { ImagePlus, PencilLine } from "lucide-react";
import type { StoryEx } from "@/modules/backend/types";
import { calculateReadingTime } from "@/lib/reading-time";

export type StoryMetadataProps = {
  story: StoryEx;
  editUrl?: string;
  coverUrl?: string;
};

export function StoryMetadata(props: StoryMetadataProps) {
  const { t, i18n } = useTranslation();

  const readingTime = calculateReadingTime(props.story.content);

  // Format the published date
  const publishedDate = new Date(props.story.created_at);
  const formattedDate = publishedDate.toLocaleDateString(i18n.language, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return (
    <div className="flex flex-row gap-3 sm:gap-6 mb-8 pb-6 border-b border-border">
      <div className="flex flex-row items-center gap-1 sm:gap-2">
        <span className="text-sm font-medium text-muted-foreground">
          {t("Stories.Published on")}:
        </span>
        <time
          dateTime={props.story.created_at}
          className="text-sm text-foreground"
        >
          {formattedDate}
        </time>
      </div>

      <div className="flex flex-row items-center gap-1 sm:gap-2">
        <span className="text-sm font-medium text-muted-foreground">
          {t("Stories.Reading time")}:
        </span>
        <span className="text-sm text-foreground">
          {readingTime} {t("Stories.min read")}
        </span>
      </div>

      {(props.editUrl !== undefined || props.coverUrl !== undefined) && (
        <div className="flex-1 flex justify-end gap-4">
          {props.editUrl !== undefined && (
            <Link
              to={props.editUrl}
              className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors !no-underline hover:!no-underline"
            >
              <PencilLine className="size-3.5" />
              {t("ContentEditor.Edit Story")}
            </Link>
          )}
        </div>
      )}
    </div>
  );
}
