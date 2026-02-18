"use client";

import * as React from "react";
import { Document, Page, pdfjs } from "react-pdf";
import { ChevronLeft, ChevronRight, Download, ZoomIn, ZoomOut } from "lucide-react";
import "react-pdf/dist/Page/AnnotationLayer.css";
import "react-pdf/dist/Page/TextLayer.css";

// Configure PDF.js worker (browser-only)
pdfjs.GlobalWorkerOptions.workerSrc = new URL(
  "pdfjs-dist/build/pdf.worker.min.mjs",
  import.meta.url,
).toString();

const DEFAULT_SCALE = 1.0;
const SCALE_STEP = 0.25;
const MIN_SCALE = 0.5;
const MAX_SCALE = 3.0;

type PDFViewerInnerProps = {
  src: string;
};

function PDFViewerInner(props: PDFViewerInnerProps) {
  const [numPages, setNumPages] = React.useState<number>(0);
  const [pageNumber, setPageNumber] = React.useState(1);
  const [scale, setScale] = React.useState(DEFAULT_SCALE);
  const [containerWidth, setContainerWidth] = React.useState<number>(0);
  const containerRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (containerRef.current === null) {
      return;
    }

    const observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        setContainerWidth(entry.contentRect.width);
      }
    });

    observer.observe(containerRef.current);

    return () => observer.disconnect();
  }, []);

  function onDocumentLoadSuccess(result: { numPages: number }) {
    setNumPages(result.numPages);
    setPageNumber(1);
  }

  function goToPrevPage() {
    setPageNumber((prev) => Math.max(prev - 1, 1));
  }

  function goToNextPage() {
    setPageNumber((prev) => Math.min(prev + 1, numPages));
  }

  function zoomIn() {
    setScale((prev) => Math.min(prev + SCALE_STEP, MAX_SCALE));
  }

  function zoomOut() {
    setScale((prev) => Math.max(prev - SCALE_STEP, MIN_SCALE));
  }

  // Calculate page width based on container and scale
  const pageWidth = containerWidth > 0 ? (containerWidth - 32) * scale : undefined;

  return (
    <div ref={containerRef} className="my-4 rounded-lg border border-border overflow-hidden">
      {/* Toolbar */}
      <div className="flex items-center justify-between bg-muted/50 px-3 py-2 border-b border-border">
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={goToPrevPage}
            disabled={pageNumber <= 1}
            className="p-1.5 rounded hover:bg-muted disabled:opacity-30 disabled:cursor-not-allowed"
            aria-label="Previous page"
          >
            <ChevronLeft className="size-4" />
          </button>
          <span className="text-sm text-muted-foreground min-w-[80px] text-center">
            {pageNumber} / {numPages}
          </span>
          <button
            type="button"
            onClick={goToNextPage}
            disabled={pageNumber >= numPages}
            className="p-1.5 rounded hover:bg-muted disabled:opacity-30 disabled:cursor-not-allowed"
            aria-label="Next page"
          >
            <ChevronRight className="size-4" />
          </button>
        </div>

        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={zoomOut}
            disabled={scale <= MIN_SCALE}
            className="p-1.5 rounded hover:bg-muted disabled:opacity-30 disabled:cursor-not-allowed"
            aria-label="Zoom out"
          >
            <ZoomOut className="size-4" />
          </button>
          <span className="text-sm text-muted-foreground min-w-[50px] text-center">
            {Math.round(scale * 100)}%
          </span>
          <button
            type="button"
            onClick={zoomIn}
            disabled={scale >= MAX_SCALE}
            className="p-1.5 rounded hover:bg-muted disabled:opacity-30 disabled:cursor-not-allowed"
            aria-label="Zoom in"
          >
            <ZoomIn className="size-4" />
          </button>
          <a
            href={props.src}
            target="_blank"
            rel="noopener noreferrer"
            className="p-1.5 rounded hover:bg-muted ml-2"
            aria-label="Download PDF"
          >
            <Download className="size-4" />
          </a>
        </div>
      </div>

      {/* PDF Content */}
      <div className="overflow-auto max-h-[80vh] bg-muted/20 flex justify-center">
        <Document
          file={props.src}
          onLoadSuccess={onDocumentLoadSuccess}
          loading={
            <div className="p-8 text-center text-muted-foreground">Loading PDF...</div>
          }
          error={
            <div className="p-8 text-center text-muted-foreground">
              Failed to load PDF.{" "}
              <a href={props.src} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
                Download instead
              </a>
            </div>
          }
        >
          <Page
            pageNumber={pageNumber}
            width={pageWidth}
            renderTextLayer
            renderAnnotationLayer
          />
        </Document>
      </div>
    </div>
  );
}

export default PDFViewerInner;
