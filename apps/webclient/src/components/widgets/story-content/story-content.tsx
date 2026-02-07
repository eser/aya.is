import type { StoryEx } from "@/modules/backend/types";
import { TextContent } from "@/components/widgets/text-content";
import { StoryMetadata } from "./story-metadata";
import { StoryFooter } from "./story-footer";

export type StoryContentProps = {
  story: StoryEx;
  compiledContent: string | null;
  currentUrl: string;
  showAuthor?: boolean;
  showPublications?: boolean;
  showShare?: boolean;
  headingOffset?: number;
  editUrl?: string;
  coverUrl?: string;
};

const headingTags = ["h1", "h2", "h3", "h4", "h5", "h6"] as const;

export function StoryContent(props: StoryContentProps) {
  const {
    story,
    compiledContent,
    currentUrl,
    showAuthor = true,
    showPublications,
    showShare = true,
    headingOffset = 1,
    editUrl,
    coverUrl,
  } = props;

  const TitleTag = headingTags[headingOffset - 1] ?? "h1";

  return (
    <article className="content">
      <TitleTag>{story.title}</TitleTag>

      <StoryMetadata story={story} editUrl={editUrl} coverUrl={coverUrl} />

      <TextContent
        compiledContent={compiledContent}
        rawContent={story.content}
        shareOptions={showShare
          ? {
            title: story.title ?? "",
            summary: story.summary,
            slug: story.slug,
            currentUrl,
          }
          : undefined}
        bare
        headingOffset={headingOffset}
      />

      {showAuthor && <StoryFooter story={story} showPublications={showPublications} />}
    </article>
  );
}
