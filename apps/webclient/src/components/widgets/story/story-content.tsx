import type { StoryEx } from "@/modules/backend/types";
import { MdxContent } from "@/components/userland/mdx-content";
import { StoryMetadata } from "./story-metadata";
import { StoryShare } from "./story-share";
import { StoryFooter } from "./story-footer";

export type StoryContentProps = {
  story: StoryEx;
  compiledContent: string | null;
  currentUrl: string;
  showAuthor?: boolean;
};

export function StoryContent(props: StoryContentProps) {
  const { story, compiledContent, currentUrl, showAuthor = true } = props;

  return (
    <div className="content">
      <h1>{story.title}</h1>

      <StoryMetadata story={story} />

      <article>
        {compiledContent !== null ? (
          <MdxContent compiledSource={compiledContent} />
        ) : (
          story.content !== null && (
            <div dangerouslySetInnerHTML={{ __html: story.content }} />
          )
        )}
      </article>

      <StoryShare story={story} currentUrl={currentUrl} />

      {showAuthor && <StoryFooter story={story} />}
    </div>
  );
}
