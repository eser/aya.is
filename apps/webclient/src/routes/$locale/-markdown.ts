/**
 * Locale index domain - markdown generation utilities
 * Provides an overview of the locale with links to main sections
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/lib/markdown-middleware";
import { siteConfig, getLocaleData } from "@/config";

/**
 * Generate markdown for the locale index page
 */
export function generateLocaleIndexMarkdown(locale: string): string {
  const localeData = getLocaleData(locale);
  const localeName = localeData?.name ?? locale;

  const frontmatter = formatFrontmatter({
    title: siteConfig.name,
    locale_name: localeName,
    generated_at: new Date().toISOString(),
  });

  const content = `${siteConfig.description}

## Sections

- [Stories](/${locale}/stories.md): Articles, announcements, and content
- [News](/${locale}/news.md): Latest news and updates
- [Products](/${locale}/products.md): Open source projects and products
- [Elements](/${locale}/elements.md): Individuals and organizations
- [Events](/${locale}/events.md): Community events

## Links

- [GitHub](${siteConfig.links.github})
- [X/Twitter](${siteConfig.links.x})
- [Instagram](${siteConfig.links.instagram})
`;

  return `${frontmatter}\n\n${content}`;
}

/**
 * Register markdown handler for locale index
 * Pattern: /$locale
 */
export function registerLocaleIndexHandler(): void {
  registerMarkdownHandler("$locale", (_params, locale, _searchParams) => {
    return Promise.resolve(generateLocaleIndexMarkdown(locale));
  });
}
