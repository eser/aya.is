/**
 * News domain - markdown generation utilities
 * News is a filtered view of stories with kind="news"
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { formatDate, parseDateFromSlug } from "@/lib/date";
import { backend } from "@/modules/backend/backend";
import type { StoryEx } from "@/modules/backend/types";

/**
 * Generate markdown for the news listing page
 */
export function generateNewsListingMarkdown(
  news: StoryEx[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "News",
    generated_at: new Date().toISOString(),
  });

  const newsLinks = news.map((item) => {
    const title = item.title ?? "Untitled";
    const slug = item.slug ?? "";
    const summary = item.summary ?? "";
    const publishDate = item.slug !== null ? parseDateFromSlug(item.slug) : null;
    const dateStr = publishDate !== null ? formatDate(publishDate, locale) : "";

    let line = `- [${title}](/${locale}/news/${slug}.md)`;
    if (dateStr !== "") {
      line += ` (${dateStr})`;
    }
    if (summary !== "") {
      line += `\n  ${summary}`;
    }
    return line;
  });

  return `${frontmatter}\n\n${newsLinks.join("\n\n")}`;
}

/**
 * Register markdown handler for news listing
 * Pattern: /$locale/news
 */
export function registerNewsListingHandler(): void {
  registerMarkdownHandler("$locale/news", async (_params, locale, _searchParams) => {
    const news = await backend.getStoriesByKinds(locale, ["news"]);

    if (news === null) {
      return null;
    }

    return generateNewsListingMarkdown(news, locale);
  });
}
