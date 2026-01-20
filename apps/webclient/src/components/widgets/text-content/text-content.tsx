import { MdxContent } from "@/components/userland/mdx-content";
import { ShareOptions } from "./share-options";
import { cn } from "@/lib/utils";

export type TextContentProps = {
  title?: string;
  compiledContent: string | null;
  rawContent?: string | null;
  shareOptions?: {
    title: string;
    summary?: string | null;
    slug?: string | null;
    currentUrl: string;
  };
  bare?: boolean;
  headingOffset?: number;
  className?: string;
};

export function TextContent(props: TextContentProps) {
  const {
    title,
    compiledContent,
    rawContent = null,
    shareOptions,
    bare = false,
    headingOffset = 1,
    className,
  } = props;

  const content = (
    <>
      {title !== undefined && <h2>{title}</h2>}

      {compiledContent !== null ? (
        <MdxContent compiledSource={compiledContent} headingOffset={headingOffset} />
      ) : (
        rawContent !== null && (
          <div dangerouslySetInnerHTML={{ __html: rawContent }} />
        )
      )}

      {shareOptions !== undefined && (
        <ShareOptions
          title={shareOptions.title}
          summary={shareOptions.summary}
          content={rawContent}
          slug={shareOptions.slug}
          currentUrl={shareOptions.currentUrl}
        />
      )}
    </>
  );

  if (bare) {
    return content;
  }

  return (
    <article className={cn("content", className)}>
      {content}
    </article>
  );
}
