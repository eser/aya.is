// Canvas Renderer for Cover Generator
// Pure functions for drawing covers on HTML Canvas

import type {
  StoryData,
  CoverOptions,
  BackgroundPattern,
} from "./types.ts";
import {
  COVER_WIDTH,
  COVER_HEIGHT,
  fontFamilyMap,
} from "./types.ts";

// Scale factor for high DPI rendering
const SCALE = 2;

// Create and configure canvas
export function createCanvas(): HTMLCanvasElement {
  const canvas = document.createElement("canvas");
  canvas.width = COVER_WIDTH * SCALE;
  canvas.height = COVER_HEIGHT * SCALE;
  return canvas;
}

// Get scaled canvas context
export function getContext(canvas: HTMLCanvasElement): CanvasRenderingContext2D {
  const ctx = canvas.getContext("2d");
  if (ctx === null) {
    throw new Error("Failed to get canvas 2d context");
  }
  ctx.scale(SCALE, SCALE);
  return ctx;
}

// Main render function
export function renderCover(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  authorImage: HTMLImageElement | null,
  logoImage: HTMLImageElement | null,
  backgroundImage: HTMLImageElement | null = null,
): void {
  // Clear canvas
  ctx.clearRect(0, 0, COVER_WIDTH, COVER_HEIGHT);

  // Draw based on template
  switch (options.templateId) {
    case "classic":
      renderClassicTemplate(ctx, story, options, authorImage, logoImage, backgroundImage);
      break;
    case "bold":
      renderBoldTemplate(ctx, story, options, authorImage, logoImage, backgroundImage);
      break;
    case "minimal":
      renderMinimalTemplate(ctx, story, options, authorImage, logoImage, backgroundImage);
      break;
    case "featured":
      renderFeaturedTemplate(ctx, story, options, authorImage, logoImage, backgroundImage);
      break;
    default:
      renderClassicTemplate(ctx, story, options, authorImage, logoImage, backgroundImage);
  }
}

// Draw background with optional pattern and image
function drawBackground(
  ctx: CanvasRenderingContext2D,
  options: CoverOptions,
  backgroundImage: HTMLImageElement | null = null,
): void {
  const { backgroundColor, backgroundPattern, backgroundImageUrl, backgroundImageOpacity } = options;

  // Fill solid background color first
  ctx.fillStyle = backgroundColor;
  ctx.fillRect(0, 0, COVER_WIDTH, COVER_HEIGHT);

  // Draw background image if provided
  if (backgroundImage !== null && backgroundImageUrl !== null) {
    ctx.globalAlpha = backgroundImageOpacity / 100;

    // Calculate dimensions to cover the canvas while maintaining aspect ratio
    const imgAspect = backgroundImage.width / backgroundImage.height;
    const canvasAspect = COVER_WIDTH / COVER_HEIGHT;

    let drawWidth: number;
    let drawHeight: number;
    let drawX: number;
    let drawY: number;

    if (imgAspect > canvasAspect) {
      // Image is wider - fit height, crop width
      drawHeight = COVER_HEIGHT;
      drawWidth = COVER_HEIGHT * imgAspect;
      drawX = (COVER_WIDTH - drawWidth) / 2;
      drawY = 0;
    } else {
      // Image is taller - fit width, crop height
      drawWidth = COVER_WIDTH;
      drawHeight = COVER_WIDTH / imgAspect;
      drawX = 0;
      drawY = (COVER_HEIGHT - drawHeight) / 2;
    }

    ctx.drawImage(backgroundImage, drawX, drawY, drawWidth, drawHeight);
    ctx.globalAlpha = 1;
  }

  // Draw pattern on top
  drawPattern(ctx, backgroundPattern, options);
}

// Draw background pattern
function drawPattern(
  ctx: CanvasRenderingContext2D,
  pattern: BackgroundPattern,
  options: CoverOptions,
): void {
  if (pattern === "none") return;

  ctx.globalAlpha = 0.1;
  ctx.fillStyle = options.textColor;

  switch (pattern) {
    case "dots":
      for (let x = 20; x < COVER_WIDTH; x += 40) {
        for (let y = 20; y < COVER_HEIGHT; y += 40) {
          ctx.beginPath();
          ctx.arc(x, y, 2, 0, Math.PI * 2);
          ctx.fill();
        }
      }
      break;

    case "grid":
      ctx.strokeStyle = options.textColor;
      ctx.lineWidth = 1;
      for (let x = 0; x < COVER_WIDTH; x += 50) {
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, COVER_HEIGHT);
        ctx.stroke();
      }
      for (let y = 0; y < COVER_HEIGHT; y += 50) {
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(COVER_WIDTH, y);
        ctx.stroke();
      }
      break;

    case "diagonal":
      ctx.strokeStyle = options.textColor;
      ctx.lineWidth = 1;
      for (let i = -COVER_HEIGHT; i < COVER_WIDTH + COVER_HEIGHT; i += 30) {
        ctx.beginPath();
        ctx.moveTo(i, 0);
        ctx.lineTo(i + COVER_HEIGHT, COVER_HEIGHT);
        ctx.stroke();
      }
      break;
  }

  ctx.globalAlpha = 1;
}

// Draw rounded rectangle
function roundRect(
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  radius: number,
): void {
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + width - radius, y);
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
  ctx.lineTo(x + width, y + height - radius);
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
  ctx.lineTo(x + radius, y + height);
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();
}

// Draw logo
function drawLogo(
  ctx: CanvasRenderingContext2D,
  logoImage: HTMLImageElement | null,
  options: CoverOptions,
): void {
  if (!options.showLogo || logoImage === null) return;

  const logoSize = 40;
  const margin = options.padding;
  let x: number;
  let y: number;

  switch (options.logoPosition) {
    case "top-left":
      x = margin;
      y = margin;
      break;
    case "top-right":
      x = COVER_WIDTH - margin - logoSize;
      y = margin;
      break;
    case "bottom-left":
      x = margin;
      y = COVER_HEIGHT - margin - logoSize;
      break;
    case "bottom-right":
    default:
      x = COVER_WIDTH - margin - logoSize;
      y = COVER_HEIGHT - margin - logoSize;
      break;
  }

  ctx.globalAlpha = options.logoOpacity / 100;
  ctx.drawImage(logoImage, x, y, logoSize, logoSize);
  ctx.globalAlpha = 1;
}

// Draw author info
function drawAuthorInfo(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  authorImage: HTMLImageElement | null,
  y: number,
): void {
  if (!options.showAuthor || story.authorName === null) return;

  const avatarSize = 48;
  const startX = options.padding;

  // Draw avatar
  if (authorImage !== null) {
    ctx.save();
    ctx.beginPath();
    ctx.arc(startX + avatarSize / 2, y + avatarSize / 2, avatarSize / 2, 0, Math.PI * 2);
    ctx.closePath();
    ctx.clip();
    ctx.drawImage(authorImage, startX, y, avatarSize, avatarSize);
    ctx.restore();
  }

  // Draw name
  ctx.fillStyle = options.textColor;
  ctx.font = `500 18px ${fontFamilyMap[options.bodyFont]}`;
  ctx.textAlign = "left";
  const textX = authorImage !== null ? startX + avatarSize + 16 : startX;
  ctx.fillText(story.authorName, textX, y + avatarSize / 2 + 6);
}

// Draw date
function drawDate(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  y: number,
): void {
  if (!options.showDate || story.publishedAt === null) return;

  const date = new Date(story.publishedAt);
  const formattedDate = date.toLocaleDateString(options.locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });

  ctx.fillStyle = options.textColor;
  ctx.globalAlpha = 0.7;
  ctx.font = `400 14px ${fontFamilyMap[options.bodyFont]}`;
  ctx.textAlign = "right";
  ctx.fillText(formattedDate, COVER_WIDTH - options.padding, y);
  ctx.globalAlpha = 1;
}

// Draw story kind badge
function drawKindBadge(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  x: number,
  y: number,
): void {
  if (!options.showStoryKind) return;

  const kind = story.kindLabel;
  ctx.font = `600 12px ${fontFamilyMap[options.bodyFont]}`;
  const textWidth = ctx.measureText(kind).width;
  const paddingX = 12;
  const paddingY = 6;
  const badgeWidth = textWidth + paddingX * 2;
  const badgeHeight = 24;

  // Badge background
  ctx.fillStyle = options.accentColor;
  roundRect(ctx, x, y, badgeWidth, badgeHeight, 4);
  ctx.fill();

  // Badge text
  ctx.fillStyle = "#ffffff";
  ctx.textAlign = "left";
  ctx.fillText(kind, x + paddingX, y + badgeHeight - paddingY);
}

// Wrap text to fit width
function wrapText(
  ctx: CanvasRenderingContext2D,
  text: string,
  maxWidth: number,
): string[] {
  const words = text.split(" ");
  const lines: string[] = [];
  let currentLine = "";

  for (const word of words) {
    const testLine = currentLine.length === 0 ? word : `${currentLine} ${word}`;
    const metrics = ctx.measureText(testLine);
    if (metrics.width > maxWidth && currentLine.length > 0) {
      lines.push(currentLine);
      currentLine = word;
    } else {
      currentLine = testLine;
    }
  }
  if (currentLine.length > 0) {
    lines.push(currentLine);
  }

  return lines;
}

// === TEMPLATE RENDERERS ===

// Classic Template: Centered title, author at bottom
function renderClassicTemplate(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  authorImage: HTMLImageElement | null,
  logoImage: HTMLImageElement | null,
  backgroundImage: HTMLImageElement | null,
): void {
  drawBackground(ctx, options, backgroundImage);

  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2;

  // Draw kind badge at top
  if (options.showStoryKind) {
    drawKindBadge(ctx, story, options, padding, padding);
  }

  // Draw title (centered)
  ctx.fillStyle = options.textColor;
  const fontSize = options.titleSize;
  ctx.font = `700 ${fontSize}px ${fontFamilyMap[options.headingFont]}`;
  ctx.textAlign = "center";

  const lines = wrapText(ctx, title, maxWidth);
  const titleLineHeight = fontSize * (options.lineHeight / 100);
  const totalHeight = lines.length * titleLineHeight;
  let startY = (COVER_HEIGHT - totalHeight) / 2 + options.contentOffsetY;

  for (const line of lines) {
    ctx.fillText(line, COVER_WIDTH / 2, startY + fontSize);
    startY += titleLineHeight;
  }

  // Draw subtitle if provided
  const subtitle = options.subtitleOverride || story.summary;
  if (subtitle !== null && subtitle.length > 0) {
    startY += options.lineSpacing; // Add spacing between title and subtitle
    ctx.globalAlpha = 0.8;
    ctx.font = `400 ${options.subtitleSize}px ${fontFamilyMap[options.bodyFont]}`;
    const subtitleLines = wrapText(ctx, subtitle, maxWidth * 0.8);
    const subtitleLineHeight = options.subtitleSize * (options.lineHeight / 100);
    for (const line of subtitleLines.slice(0, 2)) {
      ctx.fillText(line, COVER_WIDTH / 2, startY + options.subtitleSize);
      startY += subtitleLineHeight;
    }
    ctx.globalAlpha = 1;
  }

  // Draw author at bottom
  drawAuthorInfo(ctx, story, options, authorImage, COVER_HEIGHT - padding - 48);

  // Draw date
  drawDate(ctx, story, options, COVER_HEIGHT - padding - 24);

  // Draw logo
  drawLogo(ctx, logoImage, options);
}

// Bold Template: Accent highlights, badges, diagonal pattern
function renderBoldTemplate(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  authorImage: HTMLImageElement | null,
  logoImage: HTMLImageElement | null,
  backgroundImage: HTMLImageElement | null,
): void {
  drawBackground(ctx, options, backgroundImage);

  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2;

  // Draw kind badge at top
  if (options.showStoryKind) {
    drawKindBadge(ctx, story, options, padding, padding);
  }

  // Draw title with accent underline
  ctx.fillStyle = options.textColor;
  const fontSize = options.titleSize;
  ctx.font = `800 ${fontSize}px ${fontFamilyMap[options.headingFont]}`;
  ctx.textAlign = "left";

  const lines = wrapText(ctx, title, maxWidth);
  const titleLineHeight = fontSize * (options.lineHeight / 100);
  let startY = COVER_HEIGHT / 2 - (lines.length * titleLineHeight) / 2 + options.contentOffsetY;

  for (const line of lines) {
    ctx.fillText(line, padding, startY + fontSize);
    startY += titleLineHeight;
  }

  // Accent underline
  ctx.fillStyle = options.accentColor;
  ctx.fillRect(padding, startY + 10, 80, 6);
  startY += 26; // Move past underline

  // Draw subtitle if provided
  const subtitle = options.subtitleOverride || story.summary;
  if (subtitle !== null && subtitle.length > 0) {
    startY += options.lineSpacing;
    ctx.fillStyle = options.textColor;
    ctx.globalAlpha = 0.8;
    ctx.font = `400 ${options.subtitleSize}px ${fontFamilyMap[options.bodyFont]}`;
    const subtitleLines = wrapText(ctx, subtitle, maxWidth);
    const subtitleLineHeight = options.subtitleSize * (options.lineHeight / 100);
    for (const line of subtitleLines.slice(0, 2)) {
      ctx.fillText(line, padding, startY + options.subtitleSize);
      startY += subtitleLineHeight;
    }
    ctx.globalAlpha = 1;
  }

  // Draw author at bottom
  drawAuthorInfo(ctx, story, options, authorImage, COVER_HEIGHT - padding - 48);

  // Draw date
  drawDate(ctx, story, options, COVER_HEIGHT - padding - 24);

  // Draw logo
  drawLogo(ctx, logoImage, options);
}

// Minimal Template: Clean, typography-focused
function renderMinimalTemplate(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  authorImage: HTMLImageElement | null,
  logoImage: HTMLImageElement | null,
  backgroundImage: HTMLImageElement | null,
): void {
  drawBackground(ctx, options, backgroundImage);

  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2;

  // Draw kind badge at top
  if (options.showStoryKind) {
    drawKindBadge(ctx, story, options, padding, padding);
  }

  // Large centered title
  ctx.fillStyle = options.textColor;
  const fontSize = options.titleSize;
  ctx.font = `700 ${fontSize}px ${fontFamilyMap[options.headingFont]}`;
  ctx.textAlign = "center";

  const lines = wrapText(ctx, title, maxWidth);
  const titleLineHeight = fontSize * (options.lineHeight / 100);
  const totalHeight = lines.length * titleLineHeight;
  let startY = (COVER_HEIGHT - totalHeight) / 2 + options.contentOffsetY;

  for (const line of lines) {
    ctx.fillText(line, COVER_WIDTH / 2, startY + fontSize);
    startY += titleLineHeight;
  }

  // Draw subtitle if provided
  const subtitle = options.subtitleOverride || story.summary;
  if (subtitle !== null && subtitle.length > 0) {
    startY += options.lineSpacing;
    ctx.globalAlpha = 0.7;
    ctx.font = `400 ${options.subtitleSize}px ${fontFamilyMap[options.bodyFont]}`;
    const subtitleLines = wrapText(ctx, subtitle, maxWidth * 0.8);
    const subtitleLineHeight = options.subtitleSize * (options.lineHeight / 100);
    for (const line of subtitleLines.slice(0, 2)) {
      ctx.fillText(line, COVER_WIDTH / 2, startY + options.subtitleSize);
      startY += subtitleLineHeight;
    }
    ctx.globalAlpha = 1;
  }

  // Draw author at bottom
  drawAuthorInfo(ctx, story, options, authorImage, COVER_HEIGHT - padding - 48);

  // Draw date
  drawDate(ctx, story, options, COVER_HEIGHT - padding - 24);

  // Draw logo
  drawLogo(ctx, logoImage, options);
}

// Featured Template: Premium look with prominent author
function renderFeaturedTemplate(
  ctx: CanvasRenderingContext2D,
  story: StoryData,
  options: CoverOptions,
  authorImage: HTMLImageElement | null,
  logoImage: HTMLImageElement | null,
  backgroundImage: HTMLImageElement | null,
): void {
  drawBackground(ctx, options, backgroundImage);

  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2 - 150; // Leave space for author

  // Draw kind badge at top
  if (options.showStoryKind) {
    drawKindBadge(ctx, story, options, padding, padding);
  }

  // Title
  ctx.fillStyle = options.textColor;
  const fontSize = options.titleSize;
  ctx.font = `700 ${fontSize}px ${fontFamilyMap[options.headingFont]}`;
  ctx.textAlign = "left";

  const lines = wrapText(ctx, title, maxWidth);
  const titleLineHeight = fontSize * (options.lineHeight / 100);
  let startY = 120 + options.contentOffsetY;

  for (const line of lines.slice(0, 3)) {
    ctx.fillText(line, padding, startY + fontSize);
    startY += titleLineHeight;
  }

  // Draw subtitle if provided
  const subtitle = options.subtitleOverride || story.summary;
  if (subtitle !== null && subtitle.length > 0) {
    startY += options.lineSpacing;
    ctx.globalAlpha = 0.8;
    ctx.font = `400 ${options.subtitleSize}px ${fontFamilyMap[options.bodyFont]}`;
    const subtitleLines = wrapText(ctx, subtitle, maxWidth);
    const subtitleLineHeight = options.subtitleSize * (options.lineHeight / 100);
    for (const line of subtitleLines.slice(0, 2)) {
      ctx.fillText(line, padding, startY + options.subtitleSize);
      startY += subtitleLineHeight;
    }
    ctx.globalAlpha = 1;
  }

  // Large author avatar on right
  if (authorImage !== null && options.showAuthor) {
    const avatarSize = 120;
    const avatarX = COVER_WIDTH - padding - avatarSize;
    const avatarY = COVER_HEIGHT / 2 - avatarSize / 2;

    ctx.save();
    ctx.beginPath();
    ctx.arc(avatarX + avatarSize / 2, avatarY + avatarSize / 2, avatarSize / 2, 0, Math.PI * 2);
    ctx.closePath();
    ctx.clip();
    ctx.drawImage(authorImage, avatarX, avatarY, avatarSize, avatarSize);
    ctx.restore();

    // Author name below avatar
    if (story.authorName !== null) {
      ctx.fillStyle = options.textColor;
      ctx.font = `600 16px ${fontFamilyMap[options.bodyFont]}`;
      ctx.textAlign = "center";
      ctx.fillText(story.authorName, avatarX + avatarSize / 2, avatarY + avatarSize + 24);
    }
  }

  // Date at bottom left
  drawDate(ctx, story, { ...options, padding }, COVER_HEIGHT - padding);

  // Draw logo
  drawLogo(ctx, logoImage, options);
}

// Export canvas as PNG blob
export function canvasToBlob(canvas: HTMLCanvasElement): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (blob === null) {
          reject(new Error("Failed to create blob from canvas"));
          return;
        }
        resolve(blob);
      },
      "image/png",
      1.0,
    );
  });
}

// Download canvas as file
export function downloadCanvas(canvas: HTMLCanvasElement, filename: string): void {
  const link = document.createElement("a");
  link.download = filename;
  link.href = canvas.toDataURL("image/png", 1.0);
  link.click();
}
