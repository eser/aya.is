// Canvas Preview Component
// Real-time rendering of the cover image

import * as React from "react";
import {
  type StoryData,
  type CoverOptions,
  COVER_WIDTH,
  COVER_HEIGHT,
} from "@/lib/cover-generator/types.ts";
import {
  createCanvas,
  getContext,
  renderCover,
} from "@/lib/cover-generator/canvas-renderer.ts";
import { loadImage, loadLogoImage } from "@/lib/cover-generator/upload.ts";
import styles from "./cover-generator.module.css";

interface CanvasPreviewProps {
  story: StoryData;
  options: CoverOptions;
  canvasRef: React.RefObject<HTMLCanvasElement | null>;
}

export function CanvasPreview(props: CanvasPreviewProps) {
  const containerRef = React.useRef<HTMLDivElement>(null);
  const internalCanvasRef = React.useRef<HTMLCanvasElement | null>(null);
  const [authorImage, setAuthorImage] = React.useState<HTMLImageElement | null>(null);
  const [logoImage, setLogoImage] = React.useState<HTMLImageElement | null>(null);
  const [backgroundImage, setBackgroundImage] = React.useState<HTMLImageElement | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [fontsLoaded, setFontsLoaded] = React.useState(false);

  // Wait for fonts to load
  React.useEffect(() => {
    const loadFonts = async () => {
      if (typeof document !== "undefined" && document.fonts !== undefined) {
        // Wait for all fonts to be ready
        await document.fonts.ready;
        // Try to load fonts explicitly if not already loaded
        const fontsToLoad = [
          '700 16px "Bree Serif"',
          '400 16px "Nunito Sans Variable"',
          '700 16px "Nunito Sans Variable"',
        ];
        for (const font of fontsToLoad) {
          if (!document.fonts.check(font)) {
            try {
              await document.fonts.load(font);
            } catch {
              // Font may not be available, continue anyway
            }
          }
        }
      }
      setFontsLoaded(true);
    };

    loadFonts();
  }, []);

  // Load images on mount
  React.useEffect(() => {
    const loadImages = async () => {
      setIsLoading(true);

      // Load author avatar if available
      if (props.story.authorAvatarUrl !== null) {
        try {
          const img = await loadImage(props.story.authorAvatarUrl);
          setAuthorImage(img);
        } catch {
          // Silently fail - will render without avatar
        }
      }

      // Load logo with text color for matching appearance
      try {
        const logo = await loadLogoImage(props.options.textColor);
        setLogoImage(logo);
      } catch {
        // Silently fail - will render without logo
      }

      setIsLoading(false);
    };

    loadImages();
  }, [props.story.authorAvatarUrl, props.options.textColor]);

  // Load background image when URL changes
  React.useEffect(() => {
    const loadBackgroundImage = async () => {
      if (props.options.backgroundImageUrl === null) {
        setBackgroundImage(null);
        return;
      }

      try {
        const img = await loadImage(props.options.backgroundImageUrl);
        setBackgroundImage(img);
      } catch {
        // Silently fail - will render without background image
        setBackgroundImage(null);
      }
    };

    loadBackgroundImage();
  }, [props.options.backgroundImageUrl]);

  // Initialize canvas
  React.useEffect(() => {
    if (internalCanvasRef.current === null) {
      const canvas = createCanvas();
      internalCanvasRef.current = canvas;

      // Update the external ref
      if (props.canvasRef.current !== canvas) {
        (props.canvasRef as React.MutableRefObject<HTMLCanvasElement | null>).current = canvas;
      }
    }
  }, [props.canvasRef]);

  // Render cover whenever options or story changes
  React.useEffect(() => {
    if (internalCanvasRef.current === null || isLoading || !fontsLoaded) return;

    const canvas = internalCanvasRef.current;
    const ctx = getContext(canvas);

    // Reset scale since getContext applies it
    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.scale(2, 2);

    renderCover(ctx, props.story, props.options, authorImage, logoImage, backgroundImage);
  }, [props.story, props.options, authorImage, logoImage, backgroundImage, isLoading, fontsLoaded]);

  // Get canvas data URL for display
  const [previewUrl, setPreviewUrl] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (internalCanvasRef.current === null || isLoading || !fontsLoaded) return;

    // Use requestAnimationFrame to ensure canvas is rendered
    const animationId = requestAnimationFrame(() => {
      if (internalCanvasRef.current !== null) {
        setPreviewUrl(internalCanvasRef.current.toDataURL("image/png"));
      }
    });

    return () => cancelAnimationFrame(animationId);
  }, [props.story, props.options, authorImage, logoImage, backgroundImage, isLoading, fontsLoaded]);

  const showLoading = isLoading || !fontsLoaded;

  return (
    <div ref={containerRef} className={styles.previewContainer}>
      {showLoading ? (
        <div className={styles.previewLoading}>
          <div className={styles.previewLoadingSpinner} />
          <span>Loading preview...</span>
        </div>
      ) : previewUrl !== null ? (
        <img
          src={previewUrl}
          alt="Cover preview"
          className={styles.previewImage}
          style={{
            aspectRatio: `${COVER_WIDTH} / ${COVER_HEIGHT}`,
          }}
        />
      ) : null}
    </div>
  );
}
