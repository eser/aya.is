// Upload Helper for Cover Generator
// Converts canvas to blob and uploads via presigned URL

import { backend } from "@/modules/backend/backend.ts";
import { canvasToBlob } from "./canvas-renderer.ts";

export interface UploadResult {
  success: boolean;
  publicUrl?: string;
  error?: string;
}

// Upload canvas as cover image to S3
export async function uploadCoverImage(
  canvas: HTMLCanvasElement,
  locale: string,
  storySlug: string,
): Promise<UploadResult> {
  try {
    // Convert canvas to blob
    const blob = await canvasToBlob(canvas);

    // Generate unique filename
    const timestamp = Date.now();
    const filename = `cover-${storySlug}-${timestamp}.png`;

    // Get presigned URL
    const presignResponse = await backend.getPresignedURL(locale, {
      filename,
      content_type: "image/png",
      purpose: "story-picture",
    });

    if (presignResponse === null) {
      return {
        success: false,
        error: "Failed to get upload URL",
      };
    }

    // Create file from blob
    const file = new File([blob], filename, { type: "image/png" });

    // Upload to S3
    const uploadSuccess = await backend.uploadToPresignedURL(
      presignResponse.upload_url,
      file,
      "image/png",
    );

    if (!uploadSuccess) {
      return {
        success: false,
        error: "Failed to upload image",
      };
    }

    return {
      success: true,
      publicUrl: presignResponse.public_url,
    };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
}

// Load image from URL (handles CORS)
export function loadImage(url: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => resolve(img);
    img.onerror = () => reject(new Error(`Failed to load image: ${url}`));
    img.src = url;
  });
}

// Load AYA logo SVG as image (simplified version for watermark)
export async function loadLogoImage(): Promise<HTMLImageElement> {
  // Simplified AYA logo icon for watermark use
  // Uses the characteristic AYA brand colors and shapes
  const svgString = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none">
      <rect x="5" y="5" width="90" height="30" rx="4" fill="white"/>
      <rect x="5" y="40" width="65" height="22" rx="4" fill="white"/>
      <rect x="5" y="67" width="90" height="28" rx="4" fill="#66CC33"/>
    </svg>
  `;

  const blob = new Blob([svgString], { type: "image/svg+xml" });
  const url = URL.createObjectURL(blob);

  try {
    const img = await loadImage(url);
    return img;
  } finally {
    URL.revokeObjectURL(url);
  }
}
