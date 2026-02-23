import * as React from "react";
import { MdxContent } from "@/components/userland/mdx-content";

type CommentMarkdownProps = {
  content: string;
  className?: string;
};

export function CommentMarkdown(props: CommentMarkdownProps) {
  const [compiledSource, setCompiledSource] = React.useState<string | null>(null);

  React.useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const { compileMdx } = await import("@/lib/mdx");
        const compiled = await compileMdx(props.content);

        if (!cancelled) {
          setCompiledSource(compiled);
        }
      } catch {
        // On compilation error, keep showing plain text fallback
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [props.content]);

  if (compiledSource === null) {
    return (
      <p className={props.className} style={{ whiteSpace: "pre-wrap" }}>
        {props.content}
      </p>
    );
  }

  return (
    <MdxContent
      compiledSource={compiledSource}
      className={props.className}
      headingOffset={3}
      includeUserlandComponents={false}
    />
  );
}
