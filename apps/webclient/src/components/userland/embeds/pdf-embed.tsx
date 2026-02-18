"use client";

import * as React from "react";
import { Download } from "lucide-react";

export type PDFEmbedProps = {
  src: string;
};

// Lazy-loaded inner component that imports react-pdf (browser-only)
const PDFViewerInner = React.lazy(() => import("./pdf-viewer-inner"));

export function PDF(props: PDFEmbedProps) {
  const [isClient, setIsClient] = React.useState(false);

  React.useEffect(() => {
    setIsClient(true);
  }, []);

  if (!isClient) {
    return (
      <div className="my-4 rounded-lg border border-border bg-muted/30 p-8 text-center">
        <a
          href={props.src}
          target="_blank"
          rel="noopener noreferrer"
          className="text-primary hover:underline inline-flex items-center gap-2"
        >
          <Download className="size-4" />
          Download PDF
        </a>
      </div>
    );
  }

  return (
    <React.Suspense
      fallback={
        <div className="my-4 rounded-lg border border-border bg-muted/30 p-8 text-center text-muted-foreground">
          Loading PDF...
        </div>
      }
    >
      <PDFViewerInner src={props.src} />
    </React.Suspense>
  );
}
