/**
 * Contents domain - markdown generation utilities
 * Contents is a filtered view of stories with kind="content"
 */
import { formatFrontmatter, formatStoryListItem } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { backend } from "@/modules/backend/backend";
import type { StoryEx } from "@/modules/backend/types";

/**
 * Generate markdown for the contents listing page
 */
export function generateContentsListingMarkdown(
  contents: StoryEx[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "Content",
    generated_at: new Date().toISOString(),
  });

  if (contents.length === 0) {
    return `${frontmatter}\n\nNo content available yet.`;
  }

  const contentLinks = contents.map((item) => formatStoryListItem(item, locale, "contents"));

  return `${frontmatter}\n\n${contentLinks.join("\n\n")}`;
}

/**
 * Register markdown handler for contents listing
 * Pattern: /$locale/contents
 */
export function registerContentsListingHandler(): void {
  registerMarkdownHandler("$locale/contents", async (_params, locale, _searchParams) => {
    const contents = await backend.getStoriesByKinds(locale, ["content"]);

    if (contents === null) {
      return null;
    }

    return generateContentsListingMarkdown(contents, locale);
  });
}
