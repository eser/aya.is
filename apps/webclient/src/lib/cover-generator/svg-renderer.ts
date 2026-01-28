// SVG Renderer for Cover Generator
// Pure functions for generating SVG covers

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

// Escape XML special characters
function escapeXml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&apos;");
}

// Escape font-family for SVG attributes (remove inner quotes)
function escapeFontFamily(fontFamily: string): string {
  // Replace double quotes with single quotes for SVG compatibility
  return fontFamily.replace(/"/g, "'");
}

// Generate pattern definitions
function generatePatternDefs(pattern: BackgroundPattern, textColor: string): string {
  if (pattern === "none") return "";

  switch (pattern) {
    case "dots":
      return `
        <pattern id="dots-pattern" x="0" y="0" width="40" height="40" patternUnits="userSpaceOnUse">
          <circle cx="20" cy="20" r="2" fill="${textColor}" opacity="0.1"/>
        </pattern>`;

    case "grid":
      return `
        <pattern id="grid-pattern" x="0" y="0" width="50" height="50" patternUnits="userSpaceOnUse">
          <path d="M 50 0 L 0 0 0 50" fill="none" stroke="${textColor}" stroke-width="1" opacity="0.1"/>
        </pattern>`;

    case "diagonal":
      return `
        <pattern id="diagonal-pattern" x="0" y="0" width="30" height="30" patternUnits="userSpaceOnUse">
          <path d="M 0 30 L 30 0" fill="none" stroke="${textColor}" stroke-width="1" opacity="0.1"/>
        </pattern>`;

    default:
      return "";
  }
}

// Generate pattern rect
function generatePatternRect(pattern: BackgroundPattern): string {
  if (pattern === "none") return "";
  return `<rect x="0" y="0" width="${COVER_WIDTH}" height="${COVER_HEIGHT}" fill="url(#${pattern}-pattern)"/>`;
}

// Wrap text for SVG (returns array of lines)
function wrapTextSvg(text: string, fontSize: number, maxWidth: number): string[] {
  // Approximate character width based on font size
  const avgCharWidth = fontSize * 0.5;
  const maxChars = Math.floor(maxWidth / avgCharWidth);

  const words = text.split(" ");
  const lines: string[] = [];
  let currentLine = "";

  for (const word of words) {
    const testLine = currentLine.length === 0 ? word : `${currentLine} ${word}`;
    if (testLine.length > maxChars && currentLine.length > 0) {
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

// Generate logo SVG
function generateLogo(options: CoverOptions): string {
  if (!options.showLogo) return "";

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

  const opacity = options.logoOpacity / 100;

  // Simplified AYA logo icon
  return `
    <g transform="translate(${x}, ${y})" opacity="${opacity}">
      <rect x="2" y="2" width="36" height="12" rx="2" fill="white"/>
      <rect x="2" y="16" width="26" height="9" rx="2" fill="white"/>
      <rect x="2" y="27" width="36" height="11" rx="2" fill="#66CC33"/>
    </g>`;
}

// Generate kind badge
function generateKindBadge(story: StoryData, options: CoverOptions, x: number, y: number): string {
  if (!options.showStoryKind) return "";

  const kind = story.kind.charAt(0).toUpperCase() + story.kind.slice(1);
  const fontSize = 12;
  const paddingX = 12;
  const paddingY = 6;
  const textWidth = kind.length * fontSize * 0.6;
  const badgeWidth = textWidth + paddingX * 2;
  const badgeHeight = 24;

  return `
    <rect x="${x}" y="${y}" width="${badgeWidth}" height="${badgeHeight}" rx="4" fill="${options.accentColor}"/>
    <text x="${x + paddingX}" y="${y + badgeHeight - paddingY}" font-family="${escapeFontFamily(fontFamilyMap[options.bodyFont])}" font-size="${fontSize}" font-weight="600" fill="#ffffff">${escapeXml(kind)}</text>`;
}

// Generate author info
function generateAuthorInfo(
  story: StoryData,
  options: CoverOptions,
  y: number,
  authorImageDataUrl: string | null,
): string {
  if (!options.showAuthor || story.authorName === null) return "";

  const avatarSize = 48;
  const startX = options.padding;
  let result = "";

  if (authorImageDataUrl !== null) {
    result += `
      <clipPath id="avatar-clip">
        <circle cx="${startX + avatarSize / 2}" cy="${y + avatarSize / 2}" r="${avatarSize / 2}"/>
      </clipPath>
      <image href="${authorImageDataUrl}" x="${startX}" y="${y}" width="${avatarSize}" height="${avatarSize}" clip-path="url(#avatar-clip)"/>`;
  }

  const textX = authorImageDataUrl !== null ? startX + avatarSize + 16 : startX;
  result += `
    <text x="${textX}" y="${y + avatarSize / 2 + 6}" font-family="${escapeFontFamily(fontFamilyMap[options.bodyFont])}" font-size="18" font-weight="500" fill="${options.textColor}">${escapeXml(story.authorName)}</text>`;

  return result;
}

// Generate date
function generateDate(story: StoryData, options: CoverOptions, y: number): string {
  if (!options.showDate || story.publishedAt === null) return "";

  const date = new Date(story.publishedAt);
  const formattedDate = date.toLocaleDateString(options.locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });

  return `
    <text x="${COVER_WIDTH - options.padding}" y="${y}" font-family="${escapeFontFamily(fontFamilyMap[options.bodyFont])}" font-size="14" font-weight="400" fill="${options.textColor}" opacity="0.7" text-anchor="end">${escapeXml(formattedDate)}</text>`;
}

// === TEMPLATE RENDERERS ===

function renderClassicTemplateSvg(
  story: StoryData,
  options: CoverOptions,
  authorImageDataUrl: string | null,
): string {
  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2;

  const baseFontSize = 56;
  const fontSize = Math.round(baseFontSize * (options.titleSize / 100));
  const lineHeight = fontSize * (options.lineSpacing / 100);

  const lines = wrapTextSvg(title, fontSize, maxWidth);
  const totalHeight = lines.length * lineHeight;
  let startY = (COVER_HEIGHT - totalHeight) / 2;

  let titleElements = "";
  for (const line of lines) {
    titleElements += `
      <text x="${COVER_WIDTH / 2}" y="${startY + fontSize}" font-family="${escapeFontFamily(fontFamilyMap[options.headingFont])}" font-size="${fontSize}" font-weight="700" fill="${options.textColor}" text-anchor="middle">${escapeXml(line)}</text>`;
    startY += lineHeight;
  }

  // Subtitle
  let subtitleElements = "";
  const subtitle = options.subtitleOverride || story.summary;
  if (subtitle !== null && subtitle.length > 0) {
    const subtitleFontSize = 20;
    const subtitleLineHeight = subtitleFontSize * (options.lineHeight / 100);
    const subtitleLines = wrapTextSvg(subtitle, subtitleFontSize, maxWidth * 0.8);
    for (const line of subtitleLines.slice(0, 2)) {
      subtitleElements += `
        <text x="${COVER_WIDTH / 2}" y="${startY + 40}" font-family="${escapeFontFamily(fontFamilyMap[options.bodyFont])}" font-size="${subtitleFontSize}" font-weight="400" fill="${options.textColor}" opacity="0.8" text-anchor="middle">${escapeXml(line)}</text>`;
      startY += subtitleLineHeight;
    }
  }

  const authorInfo = generateAuthorInfo(story, options, COVER_HEIGHT - padding - 48, authorImageDataUrl);
  const dateInfo = generateDate(story, options, COVER_HEIGHT - padding - 24);
  const logo = generateLogo(options);

  return titleElements + subtitleElements + authorInfo + dateInfo + logo;
}

function renderBoldTemplateSvg(
  story: StoryData,
  options: CoverOptions,
  authorImageDataUrl: string | null,
): string {
  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2;

  let content = "";

  // Kind badge
  if (options.showStoryKind) {
    content += generateKindBadge(story, options, padding, padding);
  }

  const baseFontSize = 52;
  const fontSize = Math.round(baseFontSize * (options.titleSize / 100));
  const lineHeight = fontSize * (options.lineSpacing / 100);

  const lines = wrapTextSvg(title, fontSize, maxWidth);
  let startY = COVER_HEIGHT / 2 - (lines.length * lineHeight) / 2;

  for (const line of lines) {
    content += `
      <text x="${padding}" y="${startY + fontSize}" font-family="${escapeFontFamily(fontFamilyMap[options.headingFont])}" font-size="${fontSize}" font-weight="800" fill="${options.textColor}">${escapeXml(line)}</text>`;
    startY += lineHeight;
  }

  // Accent underline
  content += `<rect x="${padding}" y="${startY + 10}" width="80" height="6" fill="${options.accentColor}"/>`;

  content += generateAuthorInfo(story, options, COVER_HEIGHT - padding - 48, authorImageDataUrl);
  content += generateLogo(options);

  return content;
}

function renderMinimalTemplateSvg(
  story: StoryData,
  options: CoverOptions,
): string {
  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2;

  const baseFontSize = 64;
  const fontSize = Math.round(baseFontSize * (options.titleSize / 100));
  const lineHeight = fontSize * (options.lineSpacing / 100);

  const lines = wrapTextSvg(title, fontSize, maxWidth);
  const totalHeight = lines.length * lineHeight;
  let startY = (COVER_HEIGHT - totalHeight) / 2;

  let content = "";
  for (const line of lines) {
    content += `
      <text x="${COVER_WIDTH / 2}" y="${startY + fontSize}" font-family="${escapeFontFamily(fontFamilyMap[options.headingFont])}" font-size="${fontSize}" font-weight="700" fill="${options.textColor}" text-anchor="middle">${escapeXml(line)}</text>`;
    startY += lineHeight;
  }

  content += generateLogo(options);

  return content;
}

function renderFeaturedTemplateSvg(
  story: StoryData,
  options: CoverOptions,
  authorImageDataUrl: string | null,
): string {
  const title = options.titleOverride || story.title;
  const padding = options.padding;
  const maxWidth = COVER_WIDTH - padding * 2 - 150;

  let content = "";

  // Featured badge
  const badgeText = "FEATURED";
  const badgeWidth = badgeText.length * 14 * 0.6 + 24;
  content += `
    <rect x="${padding}" y="${padding}" width="${badgeWidth}" height="30" rx="4" fill="${options.accentColor}"/>
    <text x="${padding + 12}" y="${padding + 21}" font-family="${escapeFontFamily(fontFamilyMap[options.bodyFont])}" font-size="14" font-weight="700" fill="#000000">${badgeText}</text>`;

  // Kind badge
  if (options.showStoryKind) {
    content += generateKindBadge(story, options, padding + badgeWidth + 12, padding);
  }

  const baseFontSize = 48;
  const fontSize = Math.round(baseFontSize * (options.titleSize / 100));
  const lineHeight = fontSize * (options.lineSpacing / 100);

  const lines = wrapTextSvg(title, fontSize, maxWidth);
  let startY = 120;

  for (const line of lines.slice(0, 3)) {
    content += `
      <text x="${padding}" y="${startY + fontSize}" font-family="${escapeFontFamily(fontFamilyMap[options.headingFont])}" font-size="${fontSize}" font-weight="700" fill="${options.textColor}">${escapeXml(line)}</text>`;
    startY += lineHeight;
  }

  // Large author avatar on right
  if (authorImageDataUrl !== null && options.showAuthor) {
    const avatarSize = 120;
    const avatarX = COVER_WIDTH - padding - avatarSize;
    const avatarY = COVER_HEIGHT / 2 - avatarSize / 2;

    content += `
      <clipPath id="featured-avatar-clip">
        <circle cx="${avatarX + avatarSize / 2}" cy="${avatarY + avatarSize / 2}" r="${avatarSize / 2}"/>
      </clipPath>
      <image href="${authorImageDataUrl}" x="${avatarX}" y="${avatarY}" width="${avatarSize}" height="${avatarSize}" clip-path="url(#featured-avatar-clip)"/>`;

    if (story.authorName !== null) {
      content += `
        <text x="${avatarX + avatarSize / 2}" y="${avatarY + avatarSize + 24}" font-family="${escapeFontFamily(fontFamilyMap[options.bodyFont])}" font-size="16" font-weight="600" fill="${options.textColor}" text-anchor="middle">${escapeXml(story.authorName)}</text>`;
    }
  }

  content += generateDate(story, options, COVER_HEIGHT - padding);
  content += generateLogo(options);

  return content;
}

// Main SVG render function
export function renderCoverSvg(
  story: StoryData,
  options: CoverOptions,
  authorImageDataUrl: string | null,
): string {
  const { backgroundColor, backgroundPattern, borderRadius, textColor } = options;

  // Background
  let background: string;
  if (borderRadius > 0) {
    background = `<rect x="0" y="0" width="${COVER_WIDTH}" height="${COVER_HEIGHT}" rx="${borderRadius}" fill="${backgroundColor}"/>`;
  } else {
    background = `<rect x="0" y="0" width="${COVER_WIDTH}" height="${COVER_HEIGHT}" fill="${backgroundColor}"/>`;
  }

  // Pattern
  const patternDefs = generatePatternDefs(backgroundPattern, textColor);
  const patternRect = generatePatternRect(backgroundPattern);

  // Template content
  let templateContent: string;
  switch (options.templateId) {
    case "classic":
      templateContent = renderClassicTemplateSvg(story, options, authorImageDataUrl);
      break;
    case "bold":
      templateContent = renderBoldTemplateSvg(story, options, authorImageDataUrl);
      break;
    case "minimal":
      templateContent = renderMinimalTemplateSvg(story, options);
      break;
    case "featured":
      templateContent = renderFeaturedTemplateSvg(story, options, authorImageDataUrl);
      break;
    default:
      templateContent = renderClassicTemplateSvg(story, options, authorImageDataUrl);
  }

  return `<?xml version="1.0" encoding="UTF-8"?>
<svg width="${COVER_WIDTH}" height="${COVER_HEIGHT}" viewBox="0 0 ${COVER_WIDTH} ${COVER_HEIGHT}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <defs>
    ${patternDefs}
  </defs>
  ${background}
  ${patternRect}
  ${templateContent}
</svg>`;
}

// Download SVG as file
export function downloadSvg(svgContent: string, filename: string): void {
  const blob = new Blob([svgContent], { type: "image/svg+xml" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.download = filename;
  link.href = url;
  link.click();
  URL.revokeObjectURL(url);
}
