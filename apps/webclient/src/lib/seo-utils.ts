// Pure SEO utility functions — no import.meta.env dependency
// Extracted from seo.ts for testability with Deno's test runner

export interface SeoMeta {
  title: string;
  description: string;
  url?: string;
  image?: string | null;
  type?: "website" | "article" | "profile";
  locale?: string;
  siteName?: string;
  publishedTime?: string | null;
  modifiedTime?: string | null;
  author?: string | null;
  noIndex?: boolean;
}

export interface MetaTag {
  title?: string;
  name?: string;
  property?: string;
  content?: string;
  charSet?: string;
}

export function generateMetaTags(
  meta: SeoMeta,
  defaults: { host: string; name: string; defaultImage: string },
): MetaTag[] {
  const {
    title,
    description,
    url,
    image,
    type = "website",
    locale = "en",
    siteName = defaults.name,
    publishedTime,
    modifiedTime,
    author,
    noIndex = false,
  } = meta;

  const fullTitle = title === siteName ? title : `${title} | ${siteName}`;
  const canonicalUrl = url ?? defaults.host;
  const imageUrl = image ?? defaults.defaultImage;

  const tags: MetaTag[] = [
    // Basic meta tags
    { charSet: "utf-8" },
    { name: "viewport", content: "width=device-width, initial-scale=1" },
    { title: fullTitle },
    { name: "description", content: description },

    // Open Graph (Facebook, LinkedIn, Discord, Slack)
    { property: "og:title", content: title },
    { property: "og:description", content: description },
    { property: "og:url", content: canonicalUrl },
    { property: "og:site_name", content: siteName },
    { property: "og:type", content: type },
    { property: "og:locale", content: locale.replace("-", "_") },
  ];

  // Add image only if available
  if (imageUrl !== null && imageUrl !== "") {
    tags.push(
      { property: "og:image", content: imageUrl },
      { property: "og:image:alt", content: title },
    );
  }

  // Twitter Card
  tags.push(
    { name: "twitter:card", content: imageUrl !== null && imageUrl !== "" ? "summary_large_image" : "summary" },
    { name: "twitter:title", content: title },
    { name: "twitter:description", content: description },
    { name: "twitter:url", content: canonicalUrl },
  );

  if (imageUrl !== null && imageUrl !== "") {
    tags.push({ name: "twitter:image", content: imageUrl });
  }

  // Article-specific meta tags
  if (type === "article") {
    if (publishedTime !== null && publishedTime !== undefined) {
      tags.push({ property: "article:published_time", content: publishedTime });
    }
    if (modifiedTime !== null && modifiedTime !== undefined) {
      tags.push({ property: "article:modified_time", content: modifiedTime });
    }
    if (author !== null && author !== undefined) {
      tags.push({ property: "article:author", content: author });
    }
  }

  // Robots
  if (noIndex) {
    tags.push({ name: "robots", content: "noindex, nofollow" });
  }

  return tags;
}

// Helper to build canonical URL
export function buildUrl(host: string, locale: string, ...pathSegments: string[]): string {
  const path = pathSegments.filter((s) => s !== null && s !== undefined && s !== "").join("/");
  return `${host}/${locale}${path !== "" ? `/${path}` : ""}`;
}

// Helper to generate a canonical link tag for head()
export function generateCanonicalLink(url: string): { rel: string; href: string } {
  return { rel: "canonical", href: url };
}

// Helper to compute canonical URL for stories with locale-aware fallback
export function computeStoryCanonicalUrl(
  host: string,
  story: { is_managed: boolean; properties: Record<string, unknown> | null; locale_code?: string; slug: string | null },
  viewerLocale: string,
  ...pathPrefix: string[]
): string {
  // Rule 1: managed content with external source URL
  const sourceUrl = story.properties?.source_url as string | undefined;
  if (story.is_managed === true && sourceUrl !== undefined && sourceUrl !== "") {
    return sourceUrl;
  }

  // Rule 2: fallback content → use content's original locale
  const contentLocale = story.locale_code?.trim() ?? viewerLocale;
  const effectiveLocale = contentLocale !== "" ? contentLocale : viewerLocale;

  // Always builds on main domain via host
  return buildUrl(host, effectiveLocale, ...pathPrefix, story.slug ?? "");
}

// Helper to compute Content-Language header value
export function computeContentLanguage(viewerLocale: string, contentLocale: string | undefined): string {
  const content = contentLocale?.trim() ?? viewerLocale;
  if (content !== "" && content !== viewerLocale) {
    return `${viewerLocale}, ${content}`;
  }
  return viewerLocale;
}

// Helper to truncate description to recommended length
export function truncateDescription(
  text: string | null | undefined,
  fallbackDescription: string,
  maxLength = 160,
): string {
  if (text === null || text === undefined || text === "") {
    return fallbackDescription;
  }
  if (text.length <= maxLength) {
    return text;
  }
  return `${text.slice(0, maxLength - 3)}...`;
}
