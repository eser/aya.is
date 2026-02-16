import * as React from "react";
import { createMdxComponents, runMdxSync } from "@/lib/mdx";

type MdxContentProps = {
  compiledSource: string;
  className?: string;
  headingOffset?: number;
  includeUserlandComponents?: boolean;
};

export function MdxContent(props: MdxContentProps) {
  const {
    compiledSource,
    className,
    headingOffset = 1,
    includeUserlandComponents = true,
  } = props;

  // Use useMemo to cache the MDX component and avoid re-running on every render
  const Content = React.useMemo(() => {
    try {
      return runMdxSync(compiledSource);
    } catch (error) {
      console.error("Failed to run MDX:", error);
      return null;
    }
  }, [compiledSource]);

  const components = React.useMemo(
    () => createMdxComponents(headingOffset, includeUserlandComponents),
    [headingOffset, includeUserlandComponents],
  );

  if (Content === null) {
    return (
      <div className="text-destructive">
        Error rendering content. Please try refreshing the page.
      </div>
    );
  }

  return (
    <div className={className}>
      <Content components={components} />
    </div>
  );
}
