/**
 * llms.txt generation utilities
 * Implements the llms.txt specification for LLM-friendly content discovery
 * See: https://llmstxt.org/
 */
import { siteConfig, DEFAULT_LOCALE } from "@/config";
import { backend } from "@/modules/backend/backend";

/**
 * Generate the standard llms.txt content
 * Concise site index with links to key content
 */
export async function generateLlmsTxt(): Promise<string> {
  const locale = DEFAULT_LOCALE;

  // Fetch featured profiles
  const profiles = await backend.getProfilesByKinds(locale, [
    "individual",
    "organization",
  ]);
  const products = await backend.getProfilesByKinds(locale, ["product"]);

  // Fetch recent stories
  const stories = await backend.getStoriesByKinds(locale, [
    "article",
    "announcement",
  ]);
  const featuredStories = stories?.slice(0, 5) ?? [];

  const sections: string[] = [];

  // Header
  sections.push(`# ${siteConfig.name}`);
  sections.push("");
  sections.push(`> ${siteConfig.description}`);
  sections.push("> This file provides LLM-friendly access to AYA content.");
  sections.push("");

  // Documentation section
  sections.push("## Documentation");
  sections.push("");
  sections.push(`- [About AYA](/${locale}/aya.md): Organization profile and mission`);
  sections.push(`- [Site Index](/${locale}.md): Main navigation and sections`);
  sections.push("");

  // Profiles section
  if (profiles !== null && profiles.length > 0) {
    sections.push("## Profiles");
    sections.push("");
    for (const profile of profiles.slice(0, 10)) {
      const desc = profile.description?.slice(0, 80) ?? profile.kind;
      sections.push(`- [${profile.title}](/${locale}/${profile.slug}.md): ${desc}`);
    }
    sections.push("");
  }

  // Products section
  if (products !== null && products.length > 0) {
    sections.push("## Products");
    sections.push("");
    for (const product of products.slice(0, 5)) {
      const desc = product.description?.slice(0, 80) ?? "Product";
      sections.push(`- [${product.title}](/${locale}/${product.slug}.md): ${desc}`);
    }
    sections.push("");
  }

  // Featured Stories section
  if (featuredStories.length > 0) {
    sections.push("## Featured Stories");
    sections.push("");
    for (const story of featuredStories) {
      const title = story.title ?? "Untitled";
      const slug = story.slug ?? "";
      const summary = story.summary?.slice(0, 80) ?? "";
      sections.push(`- [${title}](/${locale}/stories/${slug}.md): ${summary}`);
    }
    sections.push("");
  }

  // Optional section (for llms.txt, just links)
  sections.push("## Optional");
  sections.push("");
  sections.push(`- [All Stories](/${locale}/stories.md): Complete story listing`);
  sections.push(`- [All News](/${locale}/news.md): News archive`);
  sections.push(`- [All Products](/${locale}/products.md): Product directory`);
  sections.push(`- [All Elements](/${locale}/elements.md): Profiles directory`);
  sections.push("");

  return sections.join("\n");
}

/**
 * Generate the extended llms-full.txt content
 * Includes more content and inline summaries
 */
export async function generateLlmsFullTxt(): Promise<string> {
  const locale = DEFAULT_LOCALE;

  // Fetch all content
  const profiles = await backend.getProfilesByKinds(locale, [
    "individual",
    "organization",
  ]);
  const products = await backend.getProfilesByKinds(locale, ["product"]);
  const stories = await backend.getStoriesByKinds(locale, [
    "article",
    "announcement",
    "content",
    "presentation",
  ]);
  const news = await backend.getStoriesByKinds(locale, ["news"]);

  const sections: string[] = [];

  // Header
  sections.push(`# ${siteConfig.name} - Full Content Index`);
  sections.push("");
  sections.push(`> ${siteConfig.description}`);
  sections.push("> This file provides comprehensive LLM-friendly access to all AYA content.");
  sections.push("");

  // Documentation section
  sections.push("## Documentation");
  sections.push("");
  sections.push(`- [About AYA](/${locale}/aya.md): Organization profile and mission`);
  sections.push(`- [Site Index](/${locale}.md): Main navigation and sections`);
  sections.push("");

  // All Profiles section
  if (profiles !== null && profiles.length > 0) {
    sections.push("## All Profiles");
    sections.push("");
    for (const profile of profiles) {
      const desc = profile.description ?? "";
      sections.push(`### [${profile.title}](/${locale}/${profile.slug}.md)`);
      sections.push(`Kind: ${profile.kind}`);
      if (desc !== "") {
        sections.push(`${desc.slice(0, 200)}${desc.length > 200 ? "..." : ""}`);
      }
      sections.push("");
    }
  }

  // All Products section
  if (products !== null && products.length > 0) {
    sections.push("## All Products");
    sections.push("");
    for (const product of products) {
      const desc = product.description ?? "";
      sections.push(`### [${product.title}](/${locale}/${product.slug}.md)`);
      if (desc !== "") {
        sections.push(`${desc.slice(0, 200)}${desc.length > 200 ? "..." : ""}`);
      }
      sections.push("");
    }
  }

  // All Stories section
  if (stories !== null && stories.length > 0) {
    sections.push("## All Stories");
    sections.push("");
    for (const story of stories) {
      const title = story.title ?? "Untitled";
      const slug = story.slug ?? "";
      const summary = story.summary ?? "";
      const author = story.author_profile?.title ?? "";

      sections.push(`### [${title}](/${locale}/stories/${slug}.md)`);
      if (author !== "") {
        sections.push(`By: ${author}`);
      }
      sections.push(`Kind: ${story.kind}`);
      if (summary !== "") {
        sections.push(`${summary.slice(0, 300)}${summary.length > 300 ? "..." : ""}`);
      }
      sections.push("");
    }
  }

  // All News section
  if (news !== null && news.length > 0) {
    sections.push("## All News");
    sections.push("");
    for (const item of news) {
      const title = item.title ?? "Untitled";
      const slug = item.slug ?? "";
      const summary = item.summary ?? "";

      sections.push(`### [${title}](/${locale}/news/${slug}.md)`);
      if (summary !== "") {
        sections.push(`${summary.slice(0, 200)}${summary.length > 200 ? "..." : ""}`);
      }
      sections.push("");
    }
  }

  return sections.join("\n");
}
