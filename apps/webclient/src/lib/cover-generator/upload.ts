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

// Load AYA logo SVG as image with customizable color
export async function loadLogoImage(color: string = "#152A35"): Promise<HTMLImageElement> {
  // Official AYA logo (without background for watermark use)
  const svgString = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3000 3000" fill="none">
      <path d="M1554.25 2243.97C1527.35 2260.24 1483.62 2260.28 1456.62 2244.16L537.662 1694.99L349.594 1794.02L1456.62 2455.56C1483.62 2471.68 1527.35 2471.59 1554.25 2455.37L2102.52 2125.22L2650.79 1795.07L2462.98 1696.76L1554.25 2243.97Z" fill="${color}" fill-opacity="0.2"/>
      <path d="M535.827 1759.24L246.253 1911.67C220.79 1925.05 205.353 1941.46 204.956 1955.52C204.559 1969.62 219.052 1986.8 243.672 2001.53L1486.84 2744.45C1489.17 2745.84 1495.63 2748.19 1505.26 2748.19C1514.94 2748.19 1521.49 2745.79 1523.87 2744.4L2756.42 2002.15C2781.04 1987.37 2795.43 1970.05 2794.99 1955.95C2794.59 1941.79 2779.05 1925.39 2753.59 1912.05L2465.01 1761L1584.62 2291.13C1562.59 2304.43 1534.4 2311.77 1505.26 2311.77C1476.37 2311.77 1448.37 2304.52 1426.44 2291.43L535.827 1759.24ZM1505.26 2859.31C1476.37 2859.31 1448.33 2852.07 1426.38 2838.97L183.216 2096.05C122.462 2059.73 88.5114 2007.38 90.0501 1952.45C91.6385 1897.51 128.518 1847.08 191.257 1814.08L539.499 1630.79L1486.84 2196.9C1489.13 2198.29 1495.63 2200.64 1505.26 2200.64C1514.99 2200.64 1521.49 2198.25 1523.87 2196.86L2460.99 1632.52L2808.34 1814.31C2871.18 1847.28 2908.15 1897.61 2909.94 1952.59C2911.63 2007.48 2877.88 2059.97 2817.22 2096.48L1584.62 2838.68C1562.59 2851.97 1534.4 2859.31 1505.26 2859.31Z" fill="${color}"/>
      <path d="M380.368 1230.26L1456.62 1873.43C1483.62 1889.6 1527.35 1889.45 1554.25 1873.28L2620.07 1231.46L1496.03 642.984L380.368 1230.26Z" fill="${color}" fill-opacity="0.2"/>
      <path d="M1496.07 706.173L246.253 1364.12C220.79 1377.46 205.353 1393.86 204.956 1407.97C204.559 1422.04 219.052 1439.25 243.672 1453.94L1486.84 2196.9C1489.13 2198.29 1495.63 2200.64 1505.26 2200.64C1514.99 2200.64 1521.49 2198.25 1523.87 2196.86L2756.42 1454.61C2781.04 1439.79 2795.43 1422.51 2794.99 1408.36C2794.59 1394.25 2779.1 1377.84 2753.59 1364.51L1496.07 706.173ZM1505.26 2311.77C1476.37 2311.77 1448.37 2304.52 1426.44 2291.43L183.216 1548.45C122.462 1512.18 88.5114 1459.84 90.0501 1404.9C91.6385 1349.97 128.518 1299.54 191.257 1266.48L1495.93 579.746L2808.34 1266.82C2871.18 1299.68 2908.15 1350.06 2909.94 1405C2911.63 1459.93 2877.88 1512.38 2817.22 1548.89L1584.62 2291.13C1562.59 2304.43 1534.4 2311.77 1505.26 2311.77Z" fill="${color}"/>
      <path d="M1456.62 1725.65C1483.62 1741.77 1527.35 1741.72 1554.25 1725.45L2786.85 983.256C2876.64 929.135 2874 845.842 2780.99 797.143L1546.56 150.904C1518.61 136.318 1473.39 136.366 1445.55 151.047L218.755 796.807C125.837 845.698 123.455 928.943 213.444 982.728L1456.62 1725.65Z" fill="${color}"/>
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
