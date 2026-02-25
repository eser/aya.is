/**
 * Profile domain - markdown generation utilities
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { backend } from "@/modules/backend/backend";
import { supportedLocales, type SupportedLocaleCode, isValidLocale } from "@/config";
import type { Profile, ProfilePage, Story } from "@/modules/backend/types";

const STORIES_PER_PAGE = 12;

interface PaginationInfo {
  offset: number;
  total: number;
  currentPage: number;
  totalPages: number;
}

/**
 * Generate markdown for a profile page
 */
export function generateProfileMarkdown(
  profile: Profile,
  locale: string,
  stories: Story[],
  pagination: PaginationInfo,
): string {
  const frontmatter = formatFrontmatter({
    title: profile.title,
    slug: profile.slug,
    kind: profile.kind,
    pronouns: profile.pronouns,
  });

  // Build the content
  const sections: string[] = [];

  // Description
  if (profile.description !== null && profile.description !== undefined) {
    sections.push(profile.description);
  }

  // Links section
  const visibleLinks = profile.links?.filter((link) => !link.is_hidden) ?? [];
  if (visibleLinks.length > 0) {
    sections.push("\n## Links\n");
    for (const link of visibleLinks) {
      if (link.uri !== null && link.uri !== undefined) {
        sections.push(`- [${link.title}](${link.uri})`);
      }
    }
  }

  // Pages section
  if (profile.pages !== undefined && profile.pages.length > 0) {
    sections.push("\n## Pages\n");
    for (const page of profile.pages) {
      sections.push(`- [${page.title}](/${locale}/${profile.slug}/${page.slug}.md)`);
    }
  }

  // Stories section
  if (stories.length > 0) {
    sections.push("\n## Stories\n");
    for (const story of stories) {
      const title = story.title ?? "Untitled";
      const slug = story.slug ?? "";
      const storyLocale = story.locale_code?.trim() ?? "";

      let line = `- [${title}](/${locale}/${profile.slug}/stories/${slug}.md)`;

      // Add locale badge when story language differs from viewer's locale
      if (storyLocale !== "" && storyLocale !== locale && isValidLocale(storyLocale)) {
        const localeData = supportedLocales[storyLocale as SupportedLocaleCode];
        line += ` [${localeData.englishName}]`;
      }

      sections.push(line);
    }
  }

  // Pagination
  if (pagination.totalPages > 1) {
    sections.push("");

    // Next page link
    if (pagination.currentPage < pagination.totalPages) {
      const nextOffset = pagination.offset + STORIES_PER_PAGE;
      sections.push(`Next Page: [Page ${pagination.currentPage + 1}](/${locale}/${profile.slug}.md?offset=${nextOffset})`);
    }

    // Previous page link
    if (pagination.currentPage > 1) {
      const prevOffset = Math.max(0, pagination.offset - STORIES_PER_PAGE);
      if (prevOffset === 0) {
        sections.push(`Previous Page: [Page ${pagination.currentPage - 1}](/${locale}/${profile.slug}.md)`);
      } else {
        sections.push(`Previous Page: [Page ${pagination.currentPage - 1}](/${locale}/${profile.slug}.md?offset=${prevOffset})`);
      }
    }

    sections.push(`Paging: ${pagination.currentPage}/${pagination.totalPages}`);
  }

  return `${frontmatter}\n\n${sections.join("\n")}`;
}

/**
 * Generate markdown for a profile custom page
 */
export function generateProfilePageMarkdown(
  page: ProfilePage,
  profile: Profile,
  _locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: page.title,
    profile: profile.title,
    profile_slug: profile.slug,
    summary: page.summary,
  });

  return `${frontmatter}\n\n${page.content}`;
}

/**
 * Shared handler logic for profile markdown
 */
async function handleProfileMarkdown(
  params: Record<string, string>,
  locale: string,
  searchParams: URLSearchParams,
): Promise<string | null> {
  const { slug } = params;

  if (slug === undefined) {
    return null;
  }

  // Skip reserved paths
  const reservedPaths = ["stories", "news", "products", "elements", "activities"];
  if (reservedPaths.includes(slug)) {
    return null;
  }

  const profile = await backend.getProfile(locale, slug);

  if (profile === null) {
    return null;
  }

  // Fetch stories for this profile
  const allStories = await backend.getProfileStories(locale, slug);
  const storiesList = allStories ?? [];

  // Handle pagination
  const offset = Number(searchParams.get("offset")) || 0;
  const total = storiesList.length;
  const totalPages = Math.ceil(total / STORIES_PER_PAGE);
  const currentPage = Math.floor(offset / STORIES_PER_PAGE) + 1;

  // Get stories for current page
  const paginatedStories = storiesList.slice(offset, offset + STORIES_PER_PAGE);

  const pagination: PaginationInfo = {
    offset,
    total,
    currentPage,
    totalPages,
  };

  return generateProfileMarkdown(profile, locale, paginatedStories, pagination);
}

/**
 * Register markdown handler for profile pages
 * Patterns: /$locale/$slug and /$locale/$slug/index
 */
export function registerProfileMarkdownHandler(): void {
  // Handle /$locale/$slug.md
  registerMarkdownHandler("$locale/$slug", handleProfileMarkdown);

  // Handle /$locale/$slug/index.md
  registerMarkdownHandler("$locale/$slug/index", handleProfileMarkdown);
}

/**
 * Register markdown handler for profile custom pages
 * Pattern: /$locale/$slug/$pageslug
 */
export function registerProfilePageMarkdownHandler(): void {
  registerMarkdownHandler("$locale/$slug/$pageslug", async (params, locale, _searchParams) => {
    const { slug, pageslug } = params;

    if (slug === undefined || pageslug === undefined) {
      return null;
    }

    // Skip reserved paths
    const reservedPaths = ["stories", "news", "products", "elements", "activities"];
    if (reservedPaths.includes(slug)) {
      return null;
    }

    // Skip if pageslug is a reserved path
    if (pageslug === "stories") {
      return null;
    }

    const profile = await backend.getProfile(locale, slug);
    if (profile === null) {
      return null;
    }

    const page = await backend.getProfilePage(locale, slug, pageslug);
    if (page === null) {
      return null;
    }

    return generateProfilePageMarkdown(page, profile, locale);
  });
}
