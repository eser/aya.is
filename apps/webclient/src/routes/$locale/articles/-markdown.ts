/**
 * Articles domain - markdown generation utilities
 * Articles is a filtered view of stories with kind="articles"
 */
import { formatFrontmatter, formatStoryListItem } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { backend } from "@/modules/backend/backend";
import type { StoryEx } from "@/modules/backend/types";

/**
 * Generate markdown for the article listing page
 */
export function generateArticlesListingMarkdown(
  articles: StoryEx[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "Articles",
    generated_at: new Date().toISOString(),
  });

  const articleLinks = articles.map((item) => formatStoryListItem(item, locale, "articles"));

  return `${frontmatter}\n\n${articleLinks.join("\n\n")}`;
}

/**
 * Register markdown handler for article listing
 * Pattern: /$locale/articles
 */
export function registerArticlesListingHandler(): void {
  registerMarkdownHandler("$locale/articles", async (_params, locale, _searchParams) => {
    const articles = await backend.getStoriesByKinds(locale, ["article"]);

    if (articles === null) {
      return null;
    }

    return generateArticlesListingMarkdown(articles, locale);
  });
}
