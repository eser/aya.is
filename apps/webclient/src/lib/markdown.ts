/**
 * Shared markdown utilities for .md routes
 * Only contains truly generic helpers used across all domains
 */

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
 * Create a Response with text/markdown content type
 */
export function createMarkdownResponse(content: string): Response {
  return new Response(content, {
    status: 200,
    headers: {
      "Content-Type": "text/markdown; charset=utf-8",
    },
  });
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
