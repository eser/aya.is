import * as React from "react";
import { createMdxComponents, runMdxSync } from "@/lib/mdx";

type MdxContentProps = {
  compiledSource: string;
  className?: string;
  headingOffset?: number;
  includeUserlandComponents?: boolean;
};

type MdxErrorBoundaryProps = {
  children: React.ReactNode;
};

type MdxErrorBoundaryState = {
  error: Error | null;
};

/**
 * Error boundary that catches React rendering errors from MDX content.
 * Prevents undefined components (e.g. `<SpotifyAlbum />`) or invalid
 * JSX from crashing the entire page.
 */
class MdxErrorBoundary extends React.Component<MdxErrorBoundaryProps, MdxErrorBoundaryState> {
  constructor(props: MdxErrorBoundaryProps) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error: Error): MdxErrorBoundaryState {
    return { error };
  }

  override componentDidUpdate(prevProps: MdxErrorBoundaryProps) {
    // Reset error state when children change (i.e. new content compiled)
    if (this.state.error !== null && prevProps.children !== this.props.children) {
      this.setState({ error: null });
    }
  }

  override render() {
    if (this.state.error !== null) {
      return (
        <div className="text-destructive text-sm">
          <strong>Content rendering error:</strong>
          <pre className="mt-1 whitespace-pre-wrap text-xs opacity-80">
            {this.state.error.message}
          </pre>
        </div>
      );
    }

    return this.props.children;
  }
}

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
    <MdxErrorBoundary>
      <div className={className}>
        <Content components={components} />
      </div>
    </MdxErrorBoundary>
  );
}
