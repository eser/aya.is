/**
 * Locale index domain - markdown generation utilities
 * Provides an overview of the locale with links to main sections
 */
import { formatFrontmatter, formatStoryListItem } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { getLocaleData, siteConfig } from "@/config";
import { backend } from "@/modules/backend/backend";

const STORY_KINDS = [
  "article",
  "news",
  "announcement",
  "status",
  "content",
  "presentation",
  "activity",
];

/**
 * Generate markdown for the locale index page
 */
export async function generateLocaleIndexMarkdown(locale: string): Promise<string> {
  const localeData = getLocaleData(locale);
  const localeName = localeData?.name ?? locale;

  const frontmatter = formatFrontmatter({
    title: siteConfig.name,
    locale_name: localeName,
    generated_at: new Date().toISOString(),
  });

  let storiesSection = "";
  try {
    const allStories = await backend.getStoriesByKinds(locale, STORY_KINDS);
    if (allStories !== null && allStories.length > 0) {
      const storyLinks = allStories.map((story) => formatStoryListItem(story, locale, "stories"));
      storiesSection = `\n\n## Latest Stories\n\n${storyLinks.join("\n\n")}`;
    }
  } catch {
    // Stories fetch failed — render page without stories
  }

  const content = `${siteConfig.description}

## Sections

- [Articles](/${locale}/articles.md): Articles and blog posts
- [Stories](/${locale}/stories.md): Announcements and content
- [News](/${locale}/news.md): Latest news and updates
- [Products](/${locale}/products.md): Open source projects and products
- [Elements](/${locale}/elements.md): Individuals and organizations
- [Activities](/${locale}/activities.md): Community activities
- [Contents](/${locale}/contents.md): Community content

## Links

- [GitHub](${siteConfig.links.github})
- [X/Twitter](${siteConfig.links.x})
- [Instagram](${siteConfig.links.instagram})${storiesSection}
`;

  return `${frontmatter}\n\n${content}`;
}

/**
 * Shared handler for locale index markdown
 */
function handleLocaleIndex(
  _params: Record<string, string>,
  locale: string,
  _searchParams: URLSearchParams,
): Promise<string | null> {
  return generateLocaleIndexMarkdown(locale);
}

/**
 * Register markdown handler for locale index
 * Patterns: /$locale and /$locale/index
 */
export function registerLocaleIndexHandler(): void {
  // Handle /$locale.md
  registerMarkdownHandler("$locale", handleLocaleIndex);

  // Handle /$locale/index.md
  registerMarkdownHandler("$locale/index", handleLocaleIndex);
}
