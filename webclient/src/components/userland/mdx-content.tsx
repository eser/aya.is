import * as React from "react";
import { runMdxSync, mdxComponents } from "@/lib/mdx";

type MdxContentProps = {
  compiledSource: string;
  className?: string;
};

export function MdxContent({ compiledSource, className }: MdxContentProps) {
  // Use useMemo to cache the MDX component and avoid re-running on every render
  const Content = React.useMemo(() => {
    try {
      return runMdxSync(compiledSource);
    } catch (error) {
      console.error("Failed to run MDX:", error);
      return null;
    }
  }, [compiledSource]);

  if (!Content) {
    return (
      <div className="text-destructive">
        Error rendering content. Please try refreshing the page.
      </div>
    );
  }

  return (
    <div className={className}>
      <Content components={mdxComponents} />
    </div>
  );
}
