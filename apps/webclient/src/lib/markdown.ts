/**
 * Shared markdown utilities for .md routes
 */
import { formatDateShort, parseDateFromSlug } from "@/lib/date";
import { isValidLocale, type SupportedLocaleCode, supportedLocales } from "@/config";

/**
 * Format data as YAML frontmatter
 */
export function formatFrontmatter(
  data: Record<string, string | number | boolean | null | undefined>,
): string {
  const lines: string[] = ["---"];

  for (const [key, value] of Object.entries(data)) {
    if (value === null || value === undefined) {
      continue;
    }

    if (typeof value === "string") {
      // Escape quotes in strings and wrap in quotes if contains special chars
      const needsQuotes = value.includes(":") || value.includes('"') ||
        value.includes("\n") || value.startsWith(" ") || value.endsWith(" ");
      if (needsQuotes) {
        const escaped = value.replace(/\\/g, "\\\\").replace(/"/g, '\\"');
        lines.push(`${key}: "${escaped}"`);
      } else {
        lines.push(`${key}: ${value}`);
      }
    } else {
      lines.push(`${key}: ${value}`);
    }
  }

  lines.push("---");
  return lines.join("\n");
}

/**
 * Format a story as a markdown list item with the standard format:
 * - [title](link) [Language] by Author (Date)
 *   Summary text
 */
export function formatStoryListItem(
  story: {
    title: string | null;
    slug: string | null;
    summary: string | null;
    locale_code?: string;
    author_profile?: { title: string } | null;
  },
  locale: string,
  basePath: string,
): string {
  const title = story.title ?? "Untitled";
  const slug = story.slug ?? "";
  const summary = story.summary ?? "";
  const author = story.author_profile?.title ?? "";
  const storyLocale = story.locale_code?.trim() ?? "";

  let line = `- [${title}](/${locale}/${basePath}/${slug}.md)`;

  // Add locale badge when story language differs from viewer's locale
  if (storyLocale !== "" && storyLocale !== locale && isValidLocale(storyLocale)) {
    const localeData = supportedLocales[storyLocale as SupportedLocaleCode];
    line += ` [${localeData.englishName}]`;
  }

  if (author !== "") {
    line += ` by ${author}`;
  }

  // Add date from slug
  const publishDate = parseDateFromSlug(story.slug);
  if (publishDate !== null) {
    line += ` (${formatDateShort(publishDate, locale)})`;
  }

  if (summary !== "") {
    line += `\n  ${summary}`;
  }

  return line;
}

/**
 * Create a Response with text/markdown content type
 */
export function createMarkdownResponse(content: string, contentLanguage?: string): Response {
  const headers: Record<string, string> = {
    "Content-Type": "text/markdown; charset=utf-8",
  };
  if (contentLanguage !== undefined) {
    headers["Content-Language"] = contentLanguage;
  }
  return new Response(content, { status: 200, headers });
}

/**
 * Create a 404 Response with markdown body
 */
export function createNotFoundResponse(): Response {
  const content = `---
title: Not Found
status: 404
---

# Not Found

The requested content could not be found.
`;

  return new Response(content, {
    status: 404,
    headers: {
      "Content-Type": "text/markdown; charset=utf-8",
    },
  });
}
