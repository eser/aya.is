// SEO Meta Tag Utilities
// Generates Open Graph, Twitter Card, and standard meta tags

import { DEFAULT_LOCALE, siteConfig } from "@/config";

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

interface MetaTag {
  title?: string;
  name?: string;
  property?: string;
  content?: string;
  charSet?: string;
}

const DEFAULT_IMAGE = `${siteConfig.host}/og-image.png`;

export function generateMetaTags(meta: SeoMeta): MetaTag[] {
  const {
    title,
    description,
    url,
    image,
    type = "website",
    locale = DEFAULT_LOCALE,
    siteName = siteConfig.name,
    publishedTime,
    modifiedTime,
    author,
    noIndex = false,
  } = meta;

  const fullTitle = title === siteName ? title : `${title} | ${siteName}`;
  const canonicalUrl = url ?? siteConfig.host;
  const imageUrl = image ?? DEFAULT_IMAGE;

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
export function buildUrl(locale: string, ...pathSegments: string[]): string {
  const path = pathSegments.filter((s) => s !== null && s !== undefined && s !== "").join("/");
  return `${siteConfig.host}/${locale}${path !== "" ? `/${path}` : ""}`;
}

// Helper to truncate description to recommended length
export function truncateDescription(text: string | null | undefined, maxLength = 160): string {
  if (text === null || text === undefined || text === "") {
    return siteConfig.description;
  }
  if (text.length <= maxLength) {
    return text;
  }
  return `${text.slice(0, maxLength - 3)}...`;
}
