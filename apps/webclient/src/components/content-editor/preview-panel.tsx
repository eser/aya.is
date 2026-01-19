import * as React from "react";
import { compileMdx } from "@/lib/mdx";
import { MdxContent } from "@/components/userland/mdx-content";
import styles from "./content-editor.module.css";

type PreviewPanelProps = {
  content: string;
  debounceMs?: number;
};

export function PreviewPanel(props: PreviewPanelProps) {
  const { content, debounceMs = 300 } = props;
  const [compiledSource, setCompiledSource] = React.useState<string | null>(
    null,
  );
  const [error, setError] = React.useState<string | null>(null);
  const [isCompiling, setIsCompiling] = React.useState(false);

  // Debounced compilation
  React.useEffect(() => {
    if (content.trim() === "") {
      setCompiledSource(null);
      setError(null);
      return;
    }

    setIsCompiling(true);
    const timer = setTimeout(async () => {
      try {
        const compiled = await compileMdx(content);
        setCompiledSource(compiled);
        setError(null);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Failed to compile MDX";
        setError(errorMessage);
        setCompiledSource(null);
      } finally {
        setIsCompiling(false);
      }
    }, debounceMs);

    return () => clearTimeout(timer);
  }, [content, debounceMs]);

  if (content.trim() === "") {
    return (
      <div className={styles.previewLoading}>
        Start writing to see the preview...
      </div>
    );
  }

  if (isCompiling && compiledSource === null) {
    return <div className={styles.previewLoading}>Compiling...</div>;
  }

  if (error !== null) {
    return (
      <div className={styles.previewError}>
        <strong>Preview Error:</strong>
        <pre className="mt-2 whitespace-pre-wrap text-xs">{error}</pre>
      </div>
    );
  }

  if (compiledSource === null) {
    return <div className={styles.previewLoading}>Loading preview...</div>;
  }

  return (
    <div className={styles.previewPanel}>
      <MdxContent compiledSource={compiledSource} headingOffset={1} />
    </div>
  );
}
