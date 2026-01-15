import * as React from "react";
import { mdxComponents, runMdxSync } from "@/lib/mdx";

type MdxContentProps = {
  compiledSource: string;
  className?: string;
};

export function MdxContent(props: MdxContentProps) {
  // Use useMemo to cache the MDX component and avoid re-running on every render
  const Content = React.useMemo(() => {
    try {
      return runMdxSync(props.compiledSource);
    } catch (error) {
      console.error("Failed to run MDX:", error);
      return null;
    }
  }, [props.compiledSource]);

  if (Content === null) {
    return (
      <div className="text-destructive">
        Error rendering content. Please try refreshing the page.
      </div>
    );
  }

  return (
    <div className={props.className}>
      <Content components={mdxComponents} />
    </div>
  );
}
